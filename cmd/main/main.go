package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Orololuwa/go-backend-boilerplate/src/config"
	"github.com/Orololuwa/go-backend-boilerplate/src/handlers"
)

const portNumber = ":8085"

var app config.AppConfig
var infoLog *log.Logger
var errorLog *log.Logger

func main (){
	err := run()
	if (err != nil){
		log.Fatal(err)
	}

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

func run() error {
	app.GoEnv = "development" //This should be gotten from the environment variables

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	return nil
}