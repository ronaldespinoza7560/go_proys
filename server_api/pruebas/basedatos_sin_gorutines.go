package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	bd "github.com/ronaldespinoza7560/go_proys/server_api/basedatos"
	"time"
)




func main() {
	
	//table := "tp_suspendidos"

	//Execute the query
	start := time.Now()
	consulta := "SELECT * FROM " + "tp_suspendidos"
	resultado, err := bd.GetJSON(consulta)
	bd.CheckErr(err)
	fmt.Println(resultado)

	consulta1 := "SELECT * FROM " + "tp_suspendidos limit 80"
	resultado1, err1 := bd.GetJSON(consulta1)
	bd.CheckErr(err1)
	fmt.Println(resultado1)

	consulta2 := "SELECT * FROM " + "tp_suspendidos limit 80"
	resultado2, err2 := bd.GetJSON(consulta2)
	bd.CheckErr(err2)
	fmt.Println(resultado2)

	consulta3 := "SELECT * FROM " + "tp_suspendidos limit 80"
	resultado3, err3 := bd.GetJSON(consulta3)
	bd.CheckErr(err3)
	fmt.Println(resultado3)

	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Println(elapsed)
	
}

