package main
import (
	"fmt"
	"time"
	fu "github.com/ronaldespinoza7560/go_proys/server_api/u2000/funciones"
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
	
	fu.Wg.Add(1)
	go genera_channel(cs)

	for i:=0;i<2;i++{
		fu.Wg.Add(1)
		go extrae_alarmas_y_los_procesa(cs)
	}
	
	
	fu.Wg.Wait()
	

	

	tiempo := time.Now()
  	elapsed := tiempo.Sub(start)
  	fmt.Println(elapsed)
}
func genera_channel(out chan []string){
	//ejecuta comandos para registrar los ne y mostrar sus alarmas
	net_elems:=[]string{"BSC01","ACT4371_Paccaritambo","MBTS_AP3693_PLAZA_AYAVIRI",
	"MBTS_TP6340_LUIS_MONTERO",
	"MBTS_TP6341_BUENOS_AIRES_SULLANA",
	"MBTS_TP6342_JOSE_MARIA_PIURA",
	"MBTS_TP6343_ESTADIO_SULLANA",
	"MBTS_TP6345_TRES_CABALLOS"}
	net_elems1:=[]string{"MBTS_TP6349_LANCONES",
	"MBTS_TP6350_LETIRA",
	"MBTS_TP6351_LA_PENITA",
	"MBTS_TP6352_AMOTAPE",
	"MBTS_TP6353_CHOCANCITO",}

	net_elems2:=[]string{"MBTS_TP6202_SAGA_FALABELLA_PIURA",
	"MBTS_TP6205_PLAYA_LOBITOS",
	"MBTS_TP6212_CURA_MORI",
	"MBTS_TP6251_SULLANA",
	"MBTS_TP6252_SULLANA_AMBEV",
	"MBTS_TP6255_CRUSETA",}
	net_elems3:=[]string{"MBTS_TP6257_SANTA_CRUZ_KM_50",
	"MBTS_TP6258_FAIQUE",
	"MBTS_TP6259_SANTO_DOMINGO_DE_CHILACO",
	"MBTS_TP6260_PUEBLO_NUEVO_COLAN",
	"MBTS_TP6261_MONTERO",
	"MBTS_TP6266_TRUCK_ECO_ACUICOLA",}
	out <- net_elems
	out <- net_elems1
	out <- net_elems2
	out <- net_elems3
	fu.Wg.Done()
	close(out)
}
func extrae_alarmas_y_los_procesa(in chan []string){
	
	//se loguea al u2000 para extraer las alarmas
	for net_elems := range in{
	
		t,err:=fu.Loguear_u2000()
		print_err(err,Error_logueo)

		for _, nelem := range net_elems {
			//fmt.Println(nelem)
			
			alarms :=  fu.Recupera_alarmas_u2000(t,nelem)
			//print_err(err1,Error_recu_alarms)
			fu.Wg.Add(1)
			go fu.Procesar_alarmas_u2000(nelem,alarms)
		
		}
		t.Write("exit\n")
	}		
	
	fu.Wg.Done()
}