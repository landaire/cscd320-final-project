package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init() {
	// If your database has a password, modify this to be <user>:<pass>@<host>/<DBName>
	if db == nil {
		var err error
		db, err = sql.Open("mysql", "root@/SalesOrders")
		if err != nil {
			Log.Fatal("An error occurred when attempting to connect to the database: ", err)
		}
	}
}
