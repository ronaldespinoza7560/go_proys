package u2000
import(
	"fmt"
	telnet "github.com/ronaldespinoza7560/go_proys/server_api/telnet"
	bd "github.com/ronaldespinoza7560/go_proys/server_api/basedatos"
	"sync"
	"strings"
	"regexp"
	"time"
	
)

//funcion que se conecta al u2000
func Loguear_u2000()(t telnet.Telnet, err error){
	defer func() {
        if r := recover(); r != nil {
            fmt.Println("Recovered in Loguear_u2000", r)
        }
    }()
	defer func() {
        if r := recover(); r != nil {
            fmt.Println("Recovered - Loguear_u2000", r)
        }
    }()
	consulta_claves:=[]string{"SELECT * FROM clave limit 1"}
	claves_u2000, err := bd.Get_datos_db(consulta_claves)
	if err!=nil{
		fmt.Println("error al conseguir la clave")

		time.Sleep(100 * time.Millisecond)
		claves_u2000, err = bd.Get_datos_db(consulta_claves)
		if err!=nil{
			claves_u2000=nil
			err=nil
			fmt.Println("error al conseguir la clave")
			return
		}
		
	}
	ip_u2000:=claves_u2000[0]["ip_u2000_wireless"].(string)
	port_u2000:=claves_u2000[0]["port_u2000_wireless"].(string)
	user_u2000:=claves_u2000[0]["user_u2000_wireless"].(string)
	clave_u2000:=claves_u2000[0]["clave_u2000_wireless"].(string)
	

	t1, err3 := telnet.Dial(ip_u2000+":"+port_u2000)
	if err3 != nil {
		fmt.Println(err)
		return
	}
	lgi_com:="LGI:OP=\""+user_u2000+"\",PWD=\""+clave_u2000+"\";"+"\r\n"
	//se logea
	t1.Read("\n")
	t1.Write(lgi_com)
	
	s, err4 := t1.Read_con_tiempo("---    END",2)
	if err4 != nil {
		fmt.Println(err4)
		return
	}
	
	fmt.Println(s[:80]+"\nSe logueo al host: "+ip_u2000+":"+port_u2000+"\ncon el usuario: "+user_u2000)
	return t1,nil
}

//recupera alarmas del u2000
func Recupera_alarmas_u2000(t telnet.Telnet,nelem string)(string){
	defer func() {
        if r := recover(); r != nil {
            fmt.Println("Recovered in Recupera_alarmas_u2000", r)
        }
    }()
		reg_ne:="REG NE:NAME="+nelem+";\r\n"
		
		t.Write(reg_ne)
		s1, err := t.Read_con_tiempo("---    END",5)
		if err != nil {
			err=nil
			fmt.Println("reg ne no responde",reg_ne)
			return ""
		}
		s1=""
		fmt.Println(s1)

		lst_almaf:="LST ALMAF:;"+"\r\n"
		t.Write(lst_almaf)
		s2, err2 := t.Read_con_tiempo("---    END",20)
		if err2 != nil {
			err2=nil
			fmt.Println("almaf no responde",reg_ne)
			return ""
		}
		//fmt.Println(s2)
	return s2
	
}

//wait group que controla las gourutines.
var Wg sync.WaitGroup
var Mux sync.Mutex

//funcion que estrae un segmento de texto contenido entre dos campos. 
func extrae_campo(campo string, campo_sig string, texto string)(string){
	arr:=strings.Split(texto,campo)
    if len(arr)==1{
        return ""
	}
	arr1:=strings.Split(arr[1],campo_sig)
    return (arr1[0])
    
}

//funcion que procesa las alarmas en crudo 
func Procesar_alarmas_u2000(netelem string, alarms string, tabla_alarmas string){
	defer func() {
        if r := recover(); r != nil {
            fmt.Println("Recovered in Procesar_alarmas_u2000", r)
        }
    }()
	tabla:=tabla_alarmas
	borrar_alarmas_celdas(tabla,netelem)
	//fmt.Println(netelem,alarms)
	alarmas:=strings.Split(alarms,"ALARM  ")
	
	//inserta las nuevas alarmas a la base de datos
	for _, alm1 :=range alarmas{
		re := regexp.MustCompile("[;, ] *")
		al:=re.Split(alm1,-1)
		//fmt.Println(al)
		if len(al) < 5 { 
            continue //si la alarma no tiene datos continua al siguiente
        }
		alarm_nro, fault,alarm_typeid:=al[0],al[2],al[4]
		//fmt.Println(alarm_nro, fault,alarm_typeid)
        alarm_typename:=extrae_campo("Alarm name  =  ","\r\n",alm1)
        alarm_rised_time:=extrae_campo("Alarm raised time  =  ","\r\n",alm1)
        network_element:=netelem
		alarm_text:="ALARM  "+alm1
		//alarm_text:="ALARM  "
		estado:="current"
		//fmt.Println(alarm_typename, alarm_rised_time,network_element,alarm_text,estado)
        hoy := time.Now().Format(time.RFC3339)
		updated_at:=hoy
		created_at:=hoy
		//fmt.Println(alarm_nro, fault,alarm_typeid,alarm_typename, alarm_rised_time,network_element,alarm_text,estado,updated_at,created_at)

		cell_name_tecno2g:=extrae_campo("AF_G=",",",alm1)
        cell_name_tecno3g:=extrae_campo("AF_U=",",",alm1)
		cell_name_tecno4g:=extrae_campo("AF_L=",",",alm1)
		cell_name:=""
		tecnologia:=""
		var cnt []string
		var cn string
		if cell_name_tecno2g != ""||cell_name_tecno3g != ""||cell_name_tecno4g != ""{
			
            if cell_name_tecno2g != ""{
                cnt:=strings.Split(cell_name_tecno2g,"_")
				cn=cnt[0]
				cell_name=cn
				if len(cn)>8{
					if string(cn[8])=="c" || string(cn[8])=="C"{
						cell_name=cn[:7]
					}
				}
                
				tecnologia="2G"
				
				//fmt.Println(cell_name,tecnologia,alarm_nro,fault,alarm_typeid,alarm_typename,alarm_rised_time,network_element,alarm_text,estado,updated_at,created_at)
				if strings.Index(alarm_nro, "++")<0 && cell_name!="" {
					insertar_datos(tabla,cell_name,tecnologia,alarm_nro,fault,alarm_typeid,alarm_typename,alarm_rised_time,network_element,alarm_text,estado,updated_at,created_at)
				}
			}   
            if cell_name_tecno3g != ""{
                cnt=strings.Split(cell_name_tecno3g,"_")
                cn=cnt[0]
				cell_name=cn
				if len(cn)>8{
					if string(cn[8])=="c" || string(cn[8])=="C"{
						cell_name=cn[:7]
					}
				}
                tecnologia="3G"
				//fmt.Println(cell_name,tecnologia,alarm_nro,fault,alarm_typeid,alarm_typename,alarm_rised_time,network_element,alarm_text,estado,updated_at,created_at)
				if strings.Index(alarm_nro, "++")<0 && cell_name!="" {
					insertar_datos(tabla,cell_name,tecnologia,alarm_nro,fault,alarm_typeid,alarm_typename,alarm_rised_time,network_element,alarm_text,estado,updated_at,created_at)
				}
			}  
            if cell_name_tecno4g != ""{
                cnt=strings.Split(cell_name_tecno4g,"_")
                cn=cnt[0]
				cell_name=cn
				if len(cn)>8{
					if string(cn[8])=="c" || string(cn[8])=="C"{
						cell_name=cn[:7]
					}
				}
                tecnologia="4G"
               // fmt.Println(cell_name,tecnologia,alarm_nro,fault,alarm_typeid,alarm_typename,alarm_rised_time,network_element,alarm_text,estado,updated_at,created_at)
			   if strings.Index(alarm_nro, "++")<0 && cell_name!="" {
					insertar_datos(tabla,cell_name,tecnologia,alarm_nro,fault,alarm_typeid,alarm_typename,alarm_rised_time,network_element,alarm_text,estado,updated_at,created_at)
				}
			}  
        }else{
				cell_name_tecnoSN:=extrae_campo("Name=",", ",alm1)
				cname:=cell_name_tecnoSN
				if cell_name_tecnoSN != ""{
					cnt=strings.Split(cname,"_")
					cn=cnt[0]
					cell_name=cn
					if len(cn)>8{
						if string(cn[8])=="c" || string(cn[8])=="C"{
							cell_name=cn[:7]
						}
					}
					
					tecnologia="2G"
					if string(cell_name[2])=="U" || string(cell_name[2])=="P"{
						tecnologia="3G"
					}
						
					if string(cell_name[2])=="L"{
						tecnologia="4G"
					}
				}
					
				//fmt.Println(cell_name,tecnologia,alarm_nro,fault,alarm_typeid,alarm_typename,alarm_rised_time,network_element,alarm_text,estado,updated_at,created_at)
				if strings.Index(alarm_nro, "++")<0 && cell_name!="" {
					insertar_datos(tabla,cell_name,tecnologia,alarm_nro,fault,alarm_typeid,alarm_typename,alarm_rised_time,network_element,alarm_text,estado,updated_at,created_at)
				}
			}
		
	}


	Wg.Done()
}

func insertar_datos(tabla,cell_name,tecnologia,alarm_nro,fault,alarm_typeid,alarm_typename,alarm_rised_time,network_element,alarm_text,estado,updated_at,created_at string){
    tabla_alarmas:=tabla
    qq_insertar:=[]string{"INSERT INTO "+tabla_alarmas +" (cell_name,tecnologia,alarm_nro,fault,alarm_typeid,alarm_typename,alarm_rised_time,network_element,alarm_text,estado,updated_at,created_at)"+
	"VALUES ('"+cell_name+"','"+tecnologia+"','"+alarm_nro+"','"+fault+"','"+alarm_typeid+"','"+alarm_typename+"','"+alarm_rised_time+"','"+network_element+"','"+alarm_text+"','"+estado+"','"+updated_at+"','"+created_at+"')"}
	Mux.Lock()
	bd.Inserta_actualiza_registros_db(qq_insertar)
	Mux.Unlock()
}

func borrar_alarmas_celdas(tabla string,ne_name string){
    tabla_alarmas:=tabla
	qq_borrar:=[]string{"DELETE FROM "+tabla_alarmas+" WHERE network_element ='"+ne_name+"'"}
	Mux.Lock()
	bd.Inserta_actualiza_registros_db(qq_borrar)
	Mux.Unlock()
}