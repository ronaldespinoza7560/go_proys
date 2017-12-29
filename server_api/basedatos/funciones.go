package basedatos

import (
	"database/sql"
	"encoding/json"
	"sync"

	"crypto/md5"
	"encoding/hex"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	//"strings";
	//bd "github.com/ronaldespinoza7560/go_proys/server_api/basedatos"
)

func checkErr(err error) (string, error) {
	if err != nil {
		return "", err
	}
	return "", nil
}

func GetJSON(sqlString string) (string, error) {
	db, err := sql.Open("mysql", Usuario+":"+Password+"@tcp("+Host+")/"+Dbname)
	checkErr(err)
	defer db.Close()
	stmt, err := db.Prepare(sqlString)
	checkErr(err)
	defer stmt.Close()

	rows, err := stmt.Query()
	checkErr(err)
	defer rows.Close()

	columns, err := rows.Columns()
	checkErr(err)

	tableData := make([]map[string]interface{}, 0)

	count := len(columns)
	values := make([]interface{}, count)
	scanArgs := make([]interface{}, count)
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		err := rows.Scan(scanArgs...)
		if err != nil {
			return "", err
		}

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
	checkErr(err)

	return string(jsonData), nil
}

func ValidarUsuario(nombre string, clave string, table string) bool {
	//var sqlString string
	hasher := md5.New()
	hasher.Write([]byte(clave))
	clave_md5 := hex.EncodeToString(hasher.Sum(nil))
	//fmt.Println(clave_md5)

	sqlString := "select count(*) as tot from " + table + " where email = '" + nombre + "' and password = '" + clave_md5 + "'"
	//fmt.Println(sqlString)

	db, err := sql.Open("mysql", Usuario+":"+Password+"@tcp("+Host+")/"+Dbname)
	checkErr(err)
	defer db.Close()

	// query
	rows, err := db.Query(sqlString)
	checkErr(err)
	
	tot := 0
	for rows.Next() {
		err = rows.Scan(&tot)
		checkErr(err)
		fmt.Println(tot)
		if tot > 0 {
			return true
		}
	}
	return false
}

func GetJSON_con_gorutines(sqlString string, Resultados map[string][]string, wg sync.WaitGroup) {
	db, err := sql.Open("mysql", Usuario+":"+Password+"@tcp("+Host+")/"+Dbname)
	checkErr(err)
	defer db.Close()
	stmt, err := db.Prepare(sqlString)
	checkErr(err)
	defer stmt.Close()

	rows, err := stmt.Query()
	checkErr(err)
	defer rows.Close()

	columns, err := rows.Columns()
	checkErr(err)

	tableData := make([]map[string]interface{}, 0)

	count := len(columns)
	values := make([]interface{}, count)
	scanArgs := make([]interface{}, count)
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		err := rows.Scan(scanArgs...)
		checkErr(err)

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
	checkErr(err)
	Resultados["consulta1"] = append(Resultados["consulta1"], string(jsonData))
	wg.Done()
	//return string(jsonData), nil
}
