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

type datos_u2000 struct{
	ip string
	port string
	user string
	clave string

}

func main() {
	start := time.Now()
	
	var u2000 datos_u2000
	
	tabla:="network_elements"  //tabla donde se almacenara las alarmas.

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
	extrae_network_elemenst_y_los_guarda_en_bd(tabla,u2000)
	
	tiempo := time.Now()
  	elapsed := tiempo.Sub(start)
  	fmt.Println(elapsed)
}


/**
* extrae network_elements del u2000 y los procesa
* recibe como entrada un arreglo de netelements
*/
func extrae_network_elemenst_y_los_guarda_en_bd(tabla string, u2000 datos_u2000){
	defer func() {
        if r := recover(); r != nil {
            fmt.Println("Recovered in extrae_alarmas_y_los_procesa", r)
        }
    }()
	//se loguea al u2000 para extraer las alarmas
	
	t,err:=fu.Loguear_u2000(u2000.ip,u2000.port,u2000.user,u2000.clave)
	print_err(err,Error_logueo)
	fu.Recupera_network_elements_u2000_y_guarda_en_db(t,tabla)
	
	t.Write("exit\n")

}