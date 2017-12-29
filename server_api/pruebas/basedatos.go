package main

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	bd "github.com/ronaldespinoza7560/go_proys/server_api/basedatos"
	"time"
)

var wg sync.WaitGroup
func main() {
	start := time.Now()
	var resultados map[string][]string
	//table := "tp_suspendidos"

	// Execute the query
	//
	// consulta := "SELECT * FROM " + "tp_suspendidos"
	// resultado, err1 := bd.GetJSON(consulta)
	// bd.CheckErr(err1)
	// fmt.Println(resultado)

	// consulta1 := "SELECT * FROM " + "tp_suspendidos limit 80"
	// resultado1, err2 := bd.GetJSON(consulta1)
	// bd.CheckErr(err2)
	// fmt.Println(resultado1)

	// t := time.Now()
	// elapsed := t.Sub(start)
	// fmt.Println(elapsed)
	wg.Add(1)
	consulta := "SELECT * FROM " + "tp_suspendidos"
	go bd.GetJSON_con_goroutines(consulta, *resultados)

	wg.Add(1)
	consulta1 := "SELECT * FROM " + "tp_suspendidos limit 80"
	go bd.GetJSON_con_goroutines(consulta, *resultados)

	wg.Wait()

	fmt.Println(resultados)
	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Println(elapsed)
}
