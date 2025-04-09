package db

import (
    "database/sql"
    _ "github.com/lib/pq"
    "log"
)

func InitDB() *sql.DB {
    connStr := "postgres://username:password@localhost/geonames?sslmode=disable"
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal(err)
    }

    err = db.Ping()
    if err != nil {
        log.Fatal(err)
    }

    return db
}