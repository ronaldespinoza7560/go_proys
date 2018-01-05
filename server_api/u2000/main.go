package main
import (
	"fmt"
	"time"
	telnet "github.com/ronaldespinoza7560/go_proys/server_api/telnet"
)

func main() {
	start := time.Now()
	t, err := telnet.Dial("172.16.108.50:31114")
	lgi_com:="LGI:OP=\"C14615\",PWD=\"R1n2ld_123\";"+"\r\n"
	if err != nil {
		fmt.Println(err)
		return
	}
	//se logea
	t.Read("\n")
	t.Write(lgi_com)
	
	s, err := t.Read_con_tiempo("---    END",2)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(s)

	//ejecuta comandos para registrar los ne y mostrar sus alarmas
	net_elems:=[]string{"BSC01","ACT4371_Paccaritambo","MBTS_AP3693_PLAZA_AYAVIRI"}
	
	for _, nelem := range net_elems {
		reg_ne:="REG NE:NAME="+nelem+";\r\n"
		
		t.Write(reg_ne)
		s1, err := t.Read_con_tiempo("---    END",2)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(s1)

		lst_almaf:="LST ALMAF:;"+"\r\n"
		t.Write(lst_almaf)
		s2, err := t.Read_con_tiempo("---    END",10)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(s2)
	
	}

	

	t.Write("exit\n")

	tiempo := time.Now()
  	elapsed := tiempo.Sub(start)
  	fmt.Println(elapsed)
}