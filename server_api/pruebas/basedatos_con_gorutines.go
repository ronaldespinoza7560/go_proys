package main

import (
	"fmt"
	"database/sql"
	"encoding/json"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	bd "github.com/ronaldespinoza7560/go_proys/server_api/basedatos"
	"time"
)

var wg sync.WaitGroup
var resultados map[string][]string
var mutex sync.Mutex



func main() {
	start := time.Now()
	resultados = make(map[string][]string)
	
	//table := "tp_suspendidos"

	// Execute the query
	// start := time.Now()
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
	go GetJSON_con_gorutines(consulta,"consulta")
	wg.Add(1)
	consulta1 := "SELECT * FROM " + "tp_suspendidos limit 80"
	go GetJSON_con_gorutines(consulta1,"consulta1")

	wg.Add(1)
	consulta2 := "SELECT * FROM " + "tp_suspendidos"
	go GetJSON_con_gorutines(consulta2,"consulta")
	wg.Add(1)
	consulta3 := "SELECT * FROM " + "tp_suspendidos limit 80"
	go GetJSON_con_gorutines(consulta3,"consulta1")

	wg.Wait()
	fmt.Println(resultados)
	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Println(elapsed)
}


func GetJSON_con_gorutines(sqlString string, conx string)(string, error) {
	db, err := sql.Open("mysql", bd.Usuario+":"+bd.Password+"@tcp("+bd.Host+")/"+bd.Dbname)
	bd.CheckErr(err)
	defer db.Close()
	stmt, err := db.Prepare(sqlString)
	bd.CheckErr(err)
	defer stmt.Close()

	rows, err := stmt.Query()
	bd.CheckErr(err)
	defer rows.Close()

	columns, err := rows.Columns()
	bd.CheckErr(err)

	tableData := make([]map[string]interface{}, 0)

	count := len(columns)
	values := make([]interface{}, count)
	scanArgs := make([]interface{}, count)
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		err := rows.Scan(scanArgs...)
		bd.CheckErr(err)

		entry := make(map[string]interface{})
		for i, col := range columns {
			v := values[i]

			b, ok := v.([]byte)
			if ok {
				entry[col] = string(b)
			} else {
				entry[col] = v
			}
		}

		tableData = append(tableData, entry)
	}

	jsonData, err := json.Marshal(tableData)
	bd.CheckErr(err)
	mutex.Lock()
	resultados[conx] = append(resultados[conx], string(jsonData))
	mutex.Unlock()
	wg.Done()
	return "", nil
}
