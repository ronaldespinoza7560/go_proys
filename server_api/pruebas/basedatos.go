package main

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	bd "github.com/ronaldespinoza7560/go_proys/server_api/basedatos"
)

func main() {

	table := "tp_suspendidos"

	// Execute the query
	consulta := "SELECT * FROM " + table
	resultado, err1 := bd.GetJSON(consulta)
	bd.CheckErr(err1)
	fmt.Println(resultado)

}
