package basedatos

import (
	"database/sql"
	"encoding/json"
	"sync"

	"crypto/md5"
	"encoding/hex"

	_ "github.com/go-sql-driver/mysql"
	//"strings";
	//bd "github.com/ronaldespinoza7560/go_proys/server_api/basedatos"
)

func checkErr(err error) ([]map[string]interface{}, error) {
	tab := make([]map[string]interface{}, 0)
	return tab, nil
}

func GetJSON(sqlString string) ([]map[string]interface{}, error) {
	
	tableData := make([]map[string]interface{}, 0)
	db, err := sql.Open("mysql", Usuario+":"+Password+"@tcp("+Host+")/"+Dbname)
	if err != nil {
		return tableData, err
	}
	defer db.Close()
	stmt, err := db.Prepare(sqlString)
	if err != nil {
		return tableData, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return tableData, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	checkErr(err)

	

	count := len(columns)
	values := make([]interface{}, count)
	scanArgs := make([]interface{}, count)
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		err := rows.Scan(scanArgs...)
		if err != nil {
			return tableData, err
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

	// jsonData, err := json.Marshal(tableData)
	// checkErr(err)

	return tableData, nil
}
/**
* valida el usuario y reponde con una structura 
* acceso: bool
* privilegios: string
*/
type privilegios struct {
	Ingreso bool
	Usuario string
	Nivel_acceso string 
	Accesos string
}

func ValidarUsuario(nom string, clave string, table string) privilegios {
	Pr:=privilegios{Ingreso:false,Usuario:nom,Nivel_acceso:"",Accesos:""}
	
	//var sqlString string
	hasher := md5.New()
	hasher.Write([]byte(clave))
	clave_md5 := hex.EncodeToString(hasher.Sum(nil))
	//fmt.Println(clave_md5)

	sqlString := "select count(*) as tot, nivel_acceso, accesos as tot from " + table + " where email = '" + nom + "' and password = '" + clave_md5 + "' limit 1"
	//fmt.Println(sqlString)

	db, err := sql.Open("mysql", Usuario+":"+Password+"@tcp("+Host+")/"+Dbname)
	if err != nil {
		return Pr
	}
	defer db.Close()

	// query
	rows, err := db.Query(sqlString)
	if err != nil {
		return Pr
	}

	tot := 0
	nivel_acceso :=""
	accesos:=""
	for rows.Next() {
		err = rows.Scan(&tot,&nivel_acceso,&accesos)
		if err != nil {
			return Pr
		}
		//	fmt.Println(tot)
		if tot > 0 {
			Pr.Ingreso=true
			Pr.Nivel_acceso=nivel_acceso
			Pr.Accesos=accesos
			return Pr
		}
	}
	
	return Pr
}


/**
*	consulta la tabla y responde datos en formato Json para ser utilizados con gorutines
*/
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
