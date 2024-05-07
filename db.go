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

func createDb() (int64, error) {
	//source env vars
	getEnvVars()

	//create db connection
	db, err := sql.Open("mysql", dsn(""))
	if err != nil {
		return 0, err
	}

	defer db.Close()

	//set db limits
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(time.Minute * 5)

	//create context for timeout
	context, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	//create database
	res, err := db.ExecContext(context, "CREATE DATABASE IF NOT EXISTS "+"url")
	if err != nil {
		log.Printf("Error occurred when creating DB\n %s. Have you set the correct credentials for the DB", err)
		return 0, err
	}

	no_rows, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	return no_rows, nil

}

// initialize datbase and load preliminary data
func initializeDb() error {
	db := openDb()
	defer db.Close()

	context, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	//create database
	_, err := db.ExecContext(context, "CREATE TABLE IF NOT EXISTS MalwareCheck(url varchar(255) NOT NULL, malware varchar(255) NOT NULL, PRIMARY KEY (url));")
	if err != nil {
		log.Printf("Error occurred when creating table\n %s", err)
		return err
	}

	_, err = db.ExecContext(context, "INSERT INTO MalwareCheck(url,malware) VALUES ('abc.com','yes'),('def.com','no'),('ghi.com','no'),('jkl.com','yes') ;")
	if err != nil {
		log.Printf("Error occurred when populating DB\n %s", err)
	}

	return nil
}

func openDb() *sql.DB {
	db, err := sql.Open("mysql", dsn("urls"))
	if err != nil {
		log.Fatal(err)
	}
	return db
}

// query db to get malware status (safe - yes or no) given url
func malwareCheck(url string) (string, error) {

	//check if no url provided
	if url == "" {
		return "", errors.New("empty url")
	}

	db := openDb()
	defer db.Close()

	context, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	//query db
	var mal string
	row, err := db.QueryContext(context, "select malware from MalwareCheck where url='"+url+"';")
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	//format response
	for row.Next() {
		err = row.Scan(&mal)
	}
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	return mal, nil

}

// data source name - format needed to access mysql server
func dsn(dbName string) string {
	var password = os.Getenv("dbPassword")
	var username = os.Getenv("dbUsername")
	var hostname = os.Getenv("hostname")
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, dbName)
}

type inputData struct {
	URL     string `json:"URL"`
	Malware string `json:"Malware"`
}

// read from entries.json and create slice of struct inputData
func readNewData() ([]inputData, error) {
	file, err := os.ReadFile("entries.json")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	var data []inputData
	err = json.Unmarshal(file, &data)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return data, nil
}

// add new entry(s) to the database based on read in data
func addNewEntry(entries []inputData) {
	db := openDb()
	defer db.Close()

	context, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	for _, val := range entries {
		//check if exists in db already
		out, _ := malwareCheck(val.URL)
		if out != "" {
			//if exists but malware state has changed - update
			if out != val.Malware {
				setMalwareState(val.URL, val.Malware)
			}
			//otherwise add to db
		} else {
			_, err := db.ExecContext(context, "INSERT INTO MalwareCheck(url,malware) VALUES ('"+val.URL+"','"+val.Malware+"');")
			if err != nil {
				log.Printf("Error occurred when populating DB\n %s", err)
			}
		}
	}
}

// alter existing db entry on malware yes or no state
func setMalwareState(url string, safe string) {
	db := openDb()
	defer db.Close()

	context, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	_, err := db.ExecContext(context, "update MalwareCheck set malware='"+safe+"' where url='"+url+"';")
	if err != nil {
		log.Fatal(err)
	}
}
