package main

import (
	"fmt"
	"github.com/chentian114/go-manager/persistent/database"
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"time"
)

func main()  {
	fmt.Println("hello go-manager")

	fmt.Println("part one connect database CRUD")
	openDB()



}

func openDB() (*sql.DB,error) {
	host := "localhost"
	dbName := "godemo"
	user := "root"
	password := "root"

	db,err := database.New(host,dbName,user,password).Set(
		database.SetCharset("utf8"),
		database.SetAllowCleartextPasswords(true),
		database.SetTimeOut(30 * time.Second),
	).Open(true)
	fmt.Println("db:",db,"err:",err)
	return db,err
}