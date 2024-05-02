package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func getEnvVars() {
	err := godotenv.Load("credentials.env")
	if err != nil {
		log.Fatal(err)
	}
}

func createDb() int64 {
	//source env vars
	getEnvVars()

	//create db connection
	db, err := sql.Open("mysql",
		dsn(""))
	if err != nil {
		log.Fatal(err)
	}

	//create context for timeout
	context, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	//create database
	res, err := db.ExecContext(context, "CREATE DATABASE IF NOT EXISTS "+"url")
	if err != nil {
		log.Printf("Error occurred when creating DB\n %s", err)
	}

	//return result of database creation
	no_rows, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	return no_rows

	// defer db.Close()
}

// load data into database
func loadData() error {
	db, err := sql.Open("mysql",
		dsn("urls"))
	if err != nil {
		log.Fatal(err)
	}
	val := db.Ping()
	return val
}

func dsn(dbName string) string {
	var password = os.Getenv("password")
	var username = os.Getenv("username")
	var hostname = os.Getenv("hostname")
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, dbName)
}
