package u2000
import(
	"fmt"
	telnet "github.com/ronaldespinoza7560/go_proys/server_api/telnet"
	bd "github.com/ronaldespinoza7560/go_proys/server_api/basedatos"
	"sync"
	"strings"
	"regexp"
	"time"
	"bytes"
	
)

//funcion que se conecta al u2000
func Loguear_u2000(ip_u2000 string,port_u2000 string,user_u2000 string,clave_u2000 string)(t telnet.Telnet, err error){
	
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

//recupera ne del u200 y los guarda en la tabla de la base de datos.
func Recupera_network_elements_u2000_y_guarda_en_db(t telnet.Telnet,tabla string){
	lst_ne:="LST NE:;\r\n"
	fmt.Println(tabla)	
	t.Write(lst_ne)
	ne1, err := t.Read_con_tiempo("---    END",60)
	if err != nil {
		err=nil
		fmt.Println("reg ne no responde",lst_ne)
	}
	
	//limpia la tabla network elements
	qq_insertar:=[]string{"DROP TABLE IF EXISTS "+tabla+"_copia;",
	"RENAME TABLE  "+tabla+" TO  "+tabla+"_copia;",
	"CREATE TABLE "+tabla+" LIKE "+tabla+"_copia;",}
	//fmt.Println(qq_insertar)
	bd.Inserta_actualiza_registros_db(qq_insertar)


	data:=[]byte(ne1)
	re := regexp.MustCompile("  +")
	replaced := re.ReplaceAll(bytes.TrimSpace(data), []byte(" "))
    
	ne2:=strings.Split(string(replaced),"\r\n")
	var ne3 []string
	re1 := regexp.MustCompile("[;, ] *")
	for _,v:=range ne2{
		
		ne3=re1.Split(strings.Trim(v," "),-1)
		if len(ne3)>2{
			switch ne3[0] {
				case
					"+++",
					"%%LST",
					"RETCODE",
					"LST",
					"NE":
					continue
				}
			ne_type:=ne3[0]
			ne_name:=ne3[1]
			ip_address:=ne3[2]
			estado:="operativo"
			hoy := time.Now().Format(time.RFC3339)
			updated_at:=hoy
			created_at:=hoy
			insertar_ne_a_db(tabla,ne_type,ne_name,ip_address,estado,updated_at,created_at)
	//		fmt.Println(tabla,ne_type,ne_name,ip_address,estado,updated_at,created_at)
		}
	}
}

func insertar_ne_a_db(tabla,ne_type,ne_name,ip_address,estado,updated_at,created_at string){
	tabla_ne:=tabla
	qq_insertar:=[]string{"INSERT INTO "+tabla_ne +" (ne_type,ne_name,ip_address,estado,updated_at,created_at)"+
	"VALUES ('"+ne_type+"','"+ne_name+"','"+ip_address+"','"+estado+"','"+updated_at+"','"+created_at+"')"}
	
	bd.Inserta_actualiza_registros_db(qq_insertar)
	
}

//recupera bts_ne del u2000
func Recupera_bts_ne_u2000(t telnet.Telnet,nelem string,ne_type string)(string){
	defer func() {
        if r := recover(); r != nil {
            fmt.Println("Recovered in Recupera_alarmas_u2000", r)
        }
    }()
		reg_ne:="REG NE:NAME="+nelem+";\r\n"
		
		t.Write(reg_ne)
		s1, err := t.Read_con_tiempo("---    END",20)
		if err != nil {
			err=nil
			fmt.Println("time out reg ne no responde",reg_ne)
			return ""
		}
		if strings.Index(s1, "NE does not Connection")>0{
			fmt.Println("NE does not Connection",reg_ne)
			return ""
		}
		if strings.Index(s1, "t Found NE")>0{
			fmt.Println("t Found NE",reg_ne)
			return ""
		}
		//fmt.Println(s1)
		out2,out3,out4,err :="","","",nil

		if strings.Index(s1, "t Found NE")<0{
			//fmt.Println(ne_type)
			
			if ne_type=="BSC6900GSMNE" || ne_type=="BSC6910GSMNE"||ne_type=="BSC6910GUNE"{
				t.Write("LST BTS:;\r\n")
				out2=leer_resultado_lst_ne(t)
			}
			
			if ne_type=="BSC6900UMTSNE"{
				t.Write("LST UNODEB:LSTFORMAT=HORIZONTAL;\r\n")
				out2=leer_resultado_lst_ne(t)
			}
			if ne_type=="BSC6910UMTSNE"{
				if nelem =="LIMRNC09"{
					t.Write("LST UNODEB:LOGICRNCID=1009,LSTFORMAT=HORIZONTAL;\r\n")
					out2=leer_resultado_lst_ne(t)
				}
                        
				if nelem =="LIMRNC10"{
					t.Write("LST UNODEB:LOGICRNCID=1010,LSTFORMAT=HORIZONTAL;\r\n")
					out2=leer_resultado_lst_ne(t)
				}
				if nelem =="LIMRNC11"{
					t.Write("LST UNODEB:LOGICRNCID=1011,LSTFORMAT=HORIZONTAL;\r\n")
					out2=leer_resultado_lst_ne(t)
				}
				if nelem =="LIMRNC12"{
					t.Write("LST UNODEB:LOGICRNCID=1012,LSTFORMAT=HORIZONTAL;\r\n")
					out2=leer_resultado_lst_ne(t)
				}
			}
			if ne_type=="BTS3900NE"{
				
				t.Write("LST NODEBFUNCTION:;\r\n")
				out2,err = t.Read_con_tiempo("---    END", 80)
				if err != nil {
					fmt.Println(err)
					out2=""
				}
				    
                t.Write("LST ENODEBFUNCTION:;\r\n")
                out3,err = t.Read_con_tiempo("---    END", 80)
				if err != nil {
					fmt.Println(err)
					out3=""
				}
				
                t.Write("LST GBTSFUNCTION:;\r\n")
                out4,err = t.Read_con_tiempo("---    END", 80)
				if err != nil {
					fmt.Println(err)
					out4=""
				}
				
				out2=out2+out3+out4
				//fmt.Println(out2)
                //extraer_btss(out5,ne,ne_type)
                
			}

		}
       
	return out2
	
}

func leer_resultado_lst_ne(t telnet.Telnet)string{
	out2:=""
	out2x:=""
	var err error =nil
	for{
		out2x:=""
		out2x,err = t.Read_con_tiempo("---    END", 80)
		//fmt.Println(out2x)
		if err != nil {
			fmt.Println(err)
		}
		if strings.Index(out2x, "To be continued...") >= 0{
			out2=out2+out2x
		}else{
			break
		}
		if strings.Index(out2x, "does not exist") >= 0{
			out2x=""
			break
		}
	} 
	out2=out2+out2x
	return out2
}

func Guardar_bts_en_db(tabla,ne_type,ne_name,site_name,updated_at,created_at,tecnologia string){
	tabla_ne:=tabla
	qq_insertar:=[]string{"INSERT INTO "+tabla_ne +" (ne_type,ne_name,site_name,updated_at,created_at,tecnologia)"+
	"VALUES ('"+ne_type+"','"+ne_name+"','"+site_name+"','"+updated_at+"','"+created_at+"','"+tecnologia+"')"}
	//fmt.Println(qq_insertar)
	bd.Inserta_actualiza_registros_db(qq_insertar)
	
}

func Limpia_la_tabla_bts_ns_de_registros_duplicados(tabla string){
	tabla_copia:=tabla+"_copia"
	qq_insertar:=[]string{
		"DROP TABLE IF EXISTS "+tabla_copia,
		"CREATE TABLE "+tabla_copia+" SELECT * FROM "+tabla+" group by ne_type,ne_name,site_name",
		"DROP TABLE IF EXISTS "+tabla,
		"CREATE TABLE "+tabla+" SELECT * FROM "+tabla_copia}
	//fmt.Println(qq_insertar)
	bd.Inserta_actualiza_registros_db(qq_insertar)
	
}
