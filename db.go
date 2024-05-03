package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
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

	defer db.Close()

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(time.Minute * 5)

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

}

// load data into database
func initializeDb() int64 {
	db := openDb()
	defer db.Close()

	context, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	//create database
	_, err := db.ExecContext(context, "CREATE TABLE IF NOT EXISTS MalwareCheck(url varchar(255) NOT NULL, malware varchar(255) NOT NULL, PRIMARY KEY (url));")
	if err != nil {
		log.Printf("Error occurred when creating table\n %s", err)
	}

	// context, cancelfunc = context.WithTimeout(context.Background(), 5*time.Second)
	//defer cancelfunc()

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

func openDb() *sql.DB {
	db, err := sql.Open("mysql", dsn("urls"))
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func malwareCheck(url string) (string, error) {

	if url == "" {
		return "", errors.New("empty url")
	}
	db := openDb()
	defer db.Close()

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
	return mal, nil

}

func dsn(dbName string) string {
	var password = os.Getenv("password")
	var username = os.Getenv("username")
	var hostname = os.Getenv("hostname")
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, dbName)
}

type inputData struct {
	URL     string `json:"URL"`
	Malware string `json:"Malware"`
}

func readNewData() []inputData {
	file, err := os.ReadFile("entries.json")
	if err != nil {
		log.Fatal(err)
	}

	var data []inputData
	err = json.Unmarshal(file, &data)

	if err != nil {
		log.Fatal(err)
	}
	return data
}

// add new entry(s) to the database based on read in data
func addNewEntry(entries []inputData) {
	db := openDb()
	defer db.Close()

	context, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	for _, val := range entries {
		// out, _ := malwareCheck(val.URL)
		// if out != ""{
		// 	if
		// }
		_, err := db.ExecContext(context, "INSERT INTO MalwareCheck(url,malware) VALUES ('"+val.URL+"','"+val.Malware+"');")
		if err != nil {
			log.Printf("Error occurred when populating DB\n %s", err)
		}
	}
}

// alter existing db entry on malware yes or no status
func setMalwareSafe(url string, safe string) {
	db := openDb()
	defer db.Close()

	context, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	_, err := db.ExecContext(context, "update MalwareCheck set malware='"+safe+"' where url='"+url+"';")
	if err != nil {
		log.Fatal(err)
	}
}
