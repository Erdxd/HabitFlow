package database

import (
	"database/sql"
	"log"
)

var db *sql.DB

func InitDb() (*sql.DB, error) {
	PsqlInfo := "host=192.168.0.12 port=5432 user=postgres password=postgres dbname=postgres sslmode=disable"
	var err error
	db, err = sql.Open("postgres", PsqlInfo)
	if err != nil {
		log.Println("Failed to connect to the database with your data")
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		log.Println("Failed to connect to the database")
		return nil, err
	}
	return db, nil

}
