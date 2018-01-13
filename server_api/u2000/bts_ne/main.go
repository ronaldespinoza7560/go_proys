package main
import (
	"fmt"
	"time"
	fu "github.com/ronaldespinoza7560/go_proys/server_api/u2000/funciones"
	bd "github.com/ronaldespinoza7560/go_proys/server_api/basedatos"
	"regexp"
)

const (
	Error_logueo="error de logueo"
	Error_recu_alarms="error al recuperar alarmas"
	Error_proc_alarms="error al procesar las alarmas"
)

func print_err(err error, msg string){
	if err != nil {
		fmt.Println(err, msg)
	}
	return
}

type datos_u2000 struct{
	ip string
	port string
	user string
	clave string

}

func main() {
	start := time.Now()
	var u2000 datos_u2000

	cs := make(chan []Ne_name_type)
	
	//query_nes:="select ne_name from network_elements where ne_name like '%mbts%' limit 400"
	query_nes:="select ne_type,ne_name from network_elements where ne_type like '%BTS3900NE%' limit 10"
	//query_nes:="select ne_name from network_elements where id > 3850"
	nro_nes:=3
	nro_de_gorutines:=3
	tabla:="bts_ne"  //tabla donde se almacenara las alarmas.

	//extrae la clave del u2000
	consulta_claves:=[]string{"SELECT * FROM clave limit 1"}
	claves_u2000, err := bd.Get_datos_db(consulta_claves)
	if err!=nil{
		fmt.Println("error al conseguir la clave")
	}
	//fmt.Println(claves_u2000)
	u2000.ip=claves_u2000[0]["ip_u2000_wireless"].(string)
	u2000.port=claves_u2000[0]["port_u2000_wireless"].(string)
	u2000.user=claves_u2000[0]["user_u2000_wireless"].(string)
	u2000.clave=claves_u2000[0]["clave_u2000_wireless"].(string)
	
	//genera el canal con nelemes agrupados en arreglos de nro_nes
	fu.Wg.Add(1)
	go genera_channel(cs, query_nes, nro_nes)

	for i:=0;i<nro_de_gorutines;i++{
		fu.Wg.Add(1)
		go extrae_alarmas_y_los_procesa(cs, tabla,u2000)
	}
		
	fu.Wg.Wait()
	
	tiempo := time.Now()
  	elapsed := tiempo.Sub(start)
  	fmt.Println(elapsed)
}

/**
* genera un canal con arreglo de network elemenst agrupados de acuerdo al nro_nes
* y recibe como entrada el query a la base de datos de netelems y la cantidad de ne elemenst que se agruparan.
*/
func genera_channel(out chan []Ne_name_type, query_ne string, nro_nes int ){
	querySelecs := []string{query_ne}
	tab, err := bd.Get_datos_db(querySelecs)
	if err !=nil{
		fmt.Println("err")
	}
	tab_n:=agrupa_nelems(tab, nro_nes)
	fmt.Println(tab_n)
	//coloca en el canal los arreglos de network elements
	for _,nelems := range tab_n{
		out <- nelems
	}

	fu.Wg.Done()
	close(out)
}

/**
* agrupa los network elements de la consulta  a la base de datos y retorna un 
*arreglo de arreglos de network elemenst
*/
type Ne_name_type struct{
	ne_name string
	ne_type string
}
func agrupa_nelems(tab []map[string]interface{}, k int)([][]Ne_name_type){

	var rsets = [][]Ne_name_type{}
	var rset = []Ne_name_type{}
	var elem Ne_name_type
	i:=1
	for _,v := range tab{
		elem.ne_name=v["ne_name"].(string)
		elem.ne_type=v["ne_type"].(string)

		rset=append(rset,elem)
		if !(i<k){
			rsets=append(rsets,rset)
			rset=nil
			i=0
		}
		i++
	}
	if len(rset)>0{
		rsets=append(rsets,rset)
	}
	return rsets
}

/**
* extrae alarmas del u2000 y los procesa
* recibe como entrada un arreglo de netelements
*/
func extrae_alarmas_y_los_procesa(in chan []Ne_name_type,tabla string, u2000 datos_u2000){
	defer fu.Wg.Done()
	defer func() {
        if r := recover(); r != nil {
            fmt.Println("Recovered in extrae_alarmas_y_los_procesa", r)
        }
    }()
	//se loguea al u2000 para extraer las alarmas
	for net_elems := range in{
		fu.Mux.Lock()
		t,err:=fu.Loguear_u2000(u2000.ip,u2000.port,u2000.user,u2000.clave)
		print_err(err,Error_logueo)
		fu.Mux.Unlock()

		for _, nelem := range net_elems {
			//fmt.Println(nelem)
			
			bts_ne_crudo :=  fu.Recupera_bts_ne_u2000(t,nelem.ne_name,nelem.ne_type)  //data cruda
			//fmt.Println(bts_ne_crudo)

			fu.Wg.Add(1)
			go Procesar_bts_ne_u2000(nelem,bts_ne_crudo,tabla)
		
		}
		fu.Mux.Lock()
		t.Write("exit\n")
		fu.Mux.Unlock()
		

	}		
	
	
}

/**
*Procesa los resultados obtenidos del u2000 al extraer los bts con sus nelems
*/

func Procesar_bts_ne_u2000(nelem Ne_name_type,data_cruda string, tabla string){
	defer fu.Wg.Done()
	//fmt.Println(nelem)
	 //fmt.Println(data_cruda)
	// fmt.Println(tabla)
	if data_cruda==""{
		return
	}

	re := regexp.MustCompile("[_;, ] *")
	btss:=re.Split(data_cruda,-1)
	//fmt.Println(btss)
	
	rem := regexp.MustCompile("[A-Z][A-Z][A-Z|0-9][0-9]+")
	
	btss0:=[]string{}
	for _,x :=range btss{
		if rem.MatchString(x)&&len(x)>5&&len(x)<11 {
			switch x {
				case
				"DBS3900","BTS3900","BTS3900B","BTS3900E":
					continue
				}
				btss0=append(btss0,x)	
		}
		
	}

	hoy := time.Now().Format(time.RFC3339)
	updated_at:=hoy
	created_at:=hoy
    for _,site_name :=range btss0{
        fmt.Println(nelem.ne_type,nelem.ne_name,site_name,updated_at,created_at)
		//guardar_bts_en_db(ne_type,ne_name,site_name,updated_at,created_at)
	}
	return
}