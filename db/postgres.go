package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

func InitDB() {
    connStr := "postgresql://" + 
        os.Getenv("PGUSER") + ":" + 
        os.Getenv("PGPASSWORD") + "@" + 
        os.Getenv("PGHOST") + "/" + 
        os.Getenv("PGDATABASE") + "?sslmode=require"

    var err error
    db, err = sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal("Ошибка подключения к БД:", err)
    }

    err = db.Ping()
    if err != nil {
        log.Fatal("БД не отвечает:", err)
    }
    
    log.Println("Успешно подключено к БД")
}

func GetDB() *sql.DB {
    return db
}