package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:root@tcp(localhost)/test?charset=utf8")
	checkErr(err)

	// insert
	// stmt, err := db.Prepare("INSERT userinfo SET username=?,departname=?,created=?")
	// checkErr(err)

	// res, err := stmt.Exec("astaxie", "rqwer", "2012-12-09")
	// checkErr(err)

	// id, err := res.LastInsertId()
	// checkErr(err)

	// fmt.Println(id)
	// update
	// stmt, err = db.Prepare("update userinfo set username=? where uid=?")
	// checkErr(err)
	// id := 10
	// res, err = stmt.Exec("jjjjjjjjjjjj", id)
	// checkErr(err)

	// affect, err := res.RowsAffected()
	// checkErr(err)

	// fmt.Println(affect)

	// query
	rows, err := db.Query("SELECT * FROM userinfo")
	checkErr(err)

	for rows.Next() {
		var uid int
		var username string
		var department string
		var created string
		err = rows.Scan(&uid, &username, &department, &created)
		checkErr(err)
		fmt.Println(uid, username, department, created)

	}

	// delete
	// stmt, err = db.Prepare("delete from userinfo where uid=?")
	// checkErr(err)

	// res, err = stmt.Exec(id)
	// checkErr(err)

	// affect, err = res.RowsAffected()
	// checkErr(err)

	// fmt.Println(affect)

	db.Close()

}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
