package basedatos

import (
	"database/sql"
	"encoding/json"
	"sync"
	//"crypto/md5"
	//"encoding/hex"
	"fmt"
	//"strings";
	//bd "github.com/ronaldespinoza7560/go_proys/server_api/basedatos"
)

func CheckErr(err error) (string, error) {
	if err != nil {
		return "", err
	}
	return "", nil
}

func GetJSON(sqlString string) (string, error) {
	db, err := sql.Open("mysql", Usuario+":"+Password+"@tcp("+Host+")/"+Dbname)
	CheckErr(err)
	defer db.Close()
	stmt, err := db.Prepare(sqlString)
	CheckErr(err)
	defer stmt.Close()

	rows, err := stmt.Query()
	CheckErr(err)
	defer rows.Close()

	columns, err := rows.Columns()
	CheckErr(err)

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
	CheckErr(err)

	return string(jsonData), nil
}

func ValidarUsuario(nombre string, clave string,table string ) (string, error) {
	var sqlString string
	//hash:=md5.Sum([]byte(clave))
	sqlString="select * from users where email = ? and password = '4a8a08f09d37b73795649038408b5f33'"
	//table, nombre, hex.EncodeToString(hash[:]))
	//fmt.Print(sqlString)

	// s := []string{"select * from ", table, " where email = '", nombre, "' and password = '",hex.EncodeToString(hash[:]),"'"};
	// sqlString=strings.Join(s, "");
	// fmt.Print(sqlString)

	db, err := sql.Open("mysql", Usuario+":"+Password+"@tcp("+Host+")/"+Dbname)
	CheckErr(err)
	defer db.Close()
	
	
	var email string // we "scan" the result in here

	rows, err1 := db.Query(sqlString,"b!") // WHERE number = 13
	CheckErr(err1)
	defer rows.Close()

	for rows.Next(){
		err = rows.Scan(&email)
		CheckErr(err)
		fmt.Print(email)
	}
	return email,nil
}



func GetJSON_con_gorutines(sqlString string, Resultados map[string][]string, wg sync.WaitGroup) {
	db, err := sql.Open("mysql", Usuario+":"+Password+"@tcp("+Host+")/"+Dbname)
	CheckErr(err)
	defer db.Close()
	stmt, err := db.Prepare(sqlString)
	CheckErr(err)
	defer stmt.Close()

	rows, err := stmt.Query()
	CheckErr(err)
	defer rows.Close()

	columns, err := rows.Columns()
	CheckErr(err)

	tableData := make([]map[string]interface{}, 0)

	count := len(columns)
	values := make([]interface{}, count)
	scanArgs := make([]interface{}, count)
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		err := rows.Scan(scanArgs...)
		CheckErr(err)

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
	CheckErr(err)
	Resultados["consulta1"] = append(Resultados["consulta1"], string(jsonData))
	wg.Done()
	//return string(jsonData), nil
}

