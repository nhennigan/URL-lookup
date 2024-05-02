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

// load env vars from credentials env file
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
	db, err := sql.Open("mysql", dsn(""))
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
func loadData() int64 {
	db, err := sql.Open("mysql", dsn("urls"))
	if err != nil {
		log.Fatal(err)
	}
	// val := db.Ping()
	// return val
	context, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	//create database
	_, err = db.ExecContext(context, "CREATE TABLE IF NOT EXISTS MalwareCheck(url varchar(255) NOT NULL, malware varchar(255) NOT NULL, PRIMARY KEY (url));")
	if err != nil {
		log.Printf("Error occurred when creating table\n %s", err)
	}

	// context, cancelfunc = context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancelfunc()

	res, err := db.ExecContext(context, "INSERT INTO MalwareCheck(url,malware) VALUES ('abc.com','yes'),('def.com','no'),('ghi.com','no'),('jkl.com','yes');")
	if err != nil {
		log.Printf("Error occurred when populating DB\n %s", err)
	}

	//return result of database population
	no_rows, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	return no_rows

}

func malwareCheck(url string) string {
	db, err := sql.Open("mysql", dsn("urls"))
	if err != nil {
		log.Fatal(err)
	}
	context, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	var mal string
	row, err := db.QueryContext(context, "select malware from MalwareCheck where url='"+url+"';")
	if err != nil {
		log.Fatal(err)
	}

	for row.Next() {
		err = row.Scan(&mal)
	}
	if err != nil {
		log.Fatal(err)
	}
	return mal

}

func dsn(dbName string) string {
	var password = os.Getenv("password")
	var username = os.Getenv("username")
	var hostname = os.Getenv("hostname")
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, dbName)
}
