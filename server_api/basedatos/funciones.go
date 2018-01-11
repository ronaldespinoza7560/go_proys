package basedatos

import (
	"fmt"
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

//retorna una tabla con los datos de la consulta realizada
func Get_datos_db(sqlString []string) ([]map[string]interface{}, error) {
	defer func() {
        if r := recover(); r != nil {
            fmt.Println("Recovered in Get_datos_db", r)
        }
    }()
	tableData := make([]map[string]interface{}, 0)
	db, err := sql.Open("mysql", Usuario+":"+Password+"@tcp("+Host+")/"+Dbname)
	if err != nil {
		fmt.Println("error al abrir mysql")
		return tableData, nil
	}
	defer db.Close()
	stmt, err1 := db.Prepare(sqlString[0])
	if err1 != nil {
		fmt.Println("error al preparar el stmt",sqlString[0])
		return tableData, nil
	}
	defer stmt.Close()

	rows, err2 := stmt.Query()
	if err2 != nil {
		fmt.Println("error al ejecutarl el query",sqlString)
		return tableData, nil
	}
	defer rows.Close()

	columns, err3 := rows.Columns()
	checkErr(err3)

	

	count := len(columns)
	values := make([]interface{}, count)
	scanArgs := make([]interface{}, count)
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		err4 := rows.Scan(scanArgs...)
		if err4 != nil {
			fmt.Println("error en rows.Scan")
			return tableData, nil
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
* Actualiza inserta un grupo de registros en la base de datos 
* recibe como parametro un arreglo de queries 
*/
func Inserta_actualiza_registros_db(sqlString []string) ([]map[string]interface{}, error) {
	defer func() {
        if r := recover(); r != nil {
            fmt.Println("Recovered in Inserta_actualiza_registros_db", r)
        }
    }()
	tableData := make([]map[string]interface{}, 0)

	db, err := sql.Open("mysql", Usuario+":"+Password+"@tcp("+Host+")/"+Dbname)
	if err != nil {
		fmt.Println("error al abrir mysql")
		return tableData, nil
	}
	defer db.Close()
	
	tx, err1 := db.Begin()
	if err1 != nil {
		fmt.Println("error al comenzar la transaccion")
		return tableData, nil
	} 

	for _, element := range sqlString {
		rows, err2 := db.Query(element)
		if err2 != nil {
			tx.Rollback()
			fmt.Println("error al ejecutar el query",element)
			return tableData, nil
		}
		defer rows.Close()
	}
	
	tx.Commit()
	
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
	defer func() {
        if r := recover(); r != nil {
            fmt.Println("Recovered in ValidarUsuario", r)
        }
    }()
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
	defer func() {
        if r := recover(); r != nil {
            fmt.Println("Recovered in GetJSON_con_gorutines", r)
        }
    }()
	db, err := sql.Open("mysql", Usuario+":"+Password+"@tcp("+Host+")/"+Dbname)
	checkErr(err)
	defer db.Close()
	stmt, err1 := db.Prepare(sqlString)
	checkErr(err1)
	defer stmt.Close()

	rows, err2 := stmt.Query()
	checkErr(err2)
	defer rows.Close()

	columns, err3 := rows.Columns()
	checkErr(err3)

	tableData := make([]map[string]interface{}, 0)

	count := len(columns)
	values := make([]interface{}, count)
	scanArgs := make([]interface{}, count)
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		err4 := rows.Scan(scanArgs...)
		checkErr(err4)

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

	jsonData, err5 := json.Marshal(tableData)
	checkErr(err5)
	Resultados["consulta1"] = append(Resultados["consulta1"], string(jsonData))
	wg.Done()
	//return string(jsonData), nil
}
