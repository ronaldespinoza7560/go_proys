package basedatos

import (
	"database/sql"
	"encoding/json"
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
