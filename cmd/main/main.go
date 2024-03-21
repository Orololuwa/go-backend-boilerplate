package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Orololuwa/go-backend-boilerplate/src/config"
	"github.com/Orololuwa/go-backend-boilerplate/src/driver"
	"github.com/Orololuwa/go-backend-boilerplate/src/handlers"
)

const portNumber = ":8085"

const dbHost = "localhost"
const dbPort = "5432"
const dbName = "bookings"
const dbUser = "orololuwa"

var app config.AppConfig
var infoLog *log.Logger
var errorLog *log.Logger

func main (){
	db, err := run()
	if (err != nil){
		log.Fatal(err)
	}
	defer db.SQL.Close()

	fmt.Println(fmt.Sprintf("Staring application on port %s", portNumber))

	srv := &http.Server{
		Addr: portNumber,
		Handler: routes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func run() (*driver.DB, error) {
	app.GoEnv = "development" //This should be gotten from the environment variables

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	// Connecto to DB
	log.Println("Connecting to dabase")
	db, err := driver.ConnectSQL(fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=", dbHost, dbPort, dbName, dbUser))
	if err != nil {
		log.Fatal("Cannot conect to database: Dying!", err)
	}
	log.Println("Connected to database")
	// 

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)

	return db, nil
}