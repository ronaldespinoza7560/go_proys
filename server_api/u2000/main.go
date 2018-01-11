package main
import (
	"fmt"
	"time"
	fu "github.com/ronaldespinoza7560/go_proys/server_api/u2000/funciones"
	bd "github.com/ronaldespinoza7560/go_proys/server_api/basedatos"
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



func main() {
	start := time.Now()

	cs := make(chan []string)
	
	//query_nes:="select ne_name from network_elements where ne_name like '%mbts%' limit 400"
	query_nes:="select ne_name from network_elements"
	//query_nes:="select ne_name from network_elements where id > 3850"
	nro_nes:=20
	nro_de_gorutines:=10
	tabla:="bts_alarmas"  //tabla donde se almacenara las alarmas.

	//genera el canal con nelemes agrupados en arreglos de nro_nes
	fu.Wg.Add(1)
	go genera_channel(cs, query_nes, nro_nes)

	for i:=0;i<nro_de_gorutines;i++{
		fu.Wg.Add(1)
		go extrae_alarmas_y_los_procesa(cs, tabla)
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
func genera_channel(out chan []string, query_ne string, nro_nes int ){
	querySelecs := []string{query_ne}
	tab, err := bd.Get_datos_db(querySelecs)
	if err !=nil{
		fmt.Println("err")
	}
	tab_n:=agrupa_nelems(tab, nro_nes)

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
func agrupa_nelems(tab []map[string]interface{}, k int)([][]string){
	var rsets =[][]string{}
	var rset= []string{}
	
	cad:=""
	i:=1
	for _,v := range tab{
		cad=v["ne_name"].(string)
		rset=append(rset,cad)
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
func extrae_alarmas_y_los_procesa(in chan []string,tabla string){
	defer func() {
        if r := recover(); r != nil {
            fmt.Println("Recovered in extrae_alarmas_y_los_procesa", r)
        }
    }()
	//se loguea al u2000 para extraer las alarmas
	for net_elems := range in{
		fu.Mux.Lock()
		t,err:=fu.Loguear_u2000()
		print_err(err,Error_logueo)
		fu.Mux.Unlock()
		for _, nelem := range net_elems {
			//fmt.Println(nelem)
			
			alarms :=  fu.Recupera_alarmas_u2000(t,nelem)
			//print_err(err1,Error_recu_alarms)
			fu.Wg.Add(1)
			go fu.Procesar_alarmas_u2000(nelem,alarms,tabla)
		
		}
		fu.Mux.Lock()
		t.Write("exit\n")
		fu.Mux.Unlock()
		

	}		
	
	fu.Wg.Done()
}