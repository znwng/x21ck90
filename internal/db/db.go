package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func Init() {
	var err error

	DB, err = sql.Open("mysql", "todo_db:Todo_db#123@tcp(127.0.0.1:3306)/todos")
	if err != nil {
		log.Fatal(err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to DB")
}

func Close() error {
	if DB == nil {
		return nil
	}

    fmt.Println("closing DB conn")
    return DB.Close()
}

