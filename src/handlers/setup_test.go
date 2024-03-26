package handlers

import (
	"log"
	"net/http"
	"os"

	"github.com/Orololuwa/go-backend-boilerplate/src/config"
	"github.com/Orololuwa/go-backend-boilerplate/src/driver"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var testApp config.AppConfig
var testDB *driver.DB


func getRoutes() http.Handler {
	testApp.GoEnv = "test" //This should be gotten from the environment variables

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	testApp.InfoLog = infoLog

	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	testApp.ErrorLog = errorLog

	// // Connecto to DB
	// log.Println("Connecting to dabase")
	// db, err := driver.ConnectSQL(fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=", dbHost, dbPort, dbName, dbUser))
	// if err != nil {
	// 	log.Fatal("Cannot conect to database: Dying!", err)
	// }
	// log.Println("Connected to database")
	// // 

	repo := NewRepo(&testApp, testDB)
	NewHandlers(repo)

	mux := chi.NewRouter()

	mux.Use(middleware.Logger)

	mux.Get("/health", Repo.Health)

	return mux;
}