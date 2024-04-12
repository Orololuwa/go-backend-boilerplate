package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Orololuwa/go-backend-boilerplate/src/config"
	"github.com/Orololuwa/go-backend-boilerplate/src/driver"
	"github.com/Orololuwa/go-backend-boilerplate/src/handlers"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

const portNumber = ":8085"

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
		Handler: routes(&app, db),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func run() (*driver.DB, error) {
	// read env files
	err := godotenv.Load(dir(".env.example"))
	if err != nil {
	  log.Fatal("Error loading .env file")
	}

	goEnv := os.Getenv("GO_ENV")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbSSL := os.Getenv("DB_SSL")
	
	// read flags
	// goEnv := flag.String("goenv", "development", "the application environment")
	// dbHost := flag.String("dbhost", "localhost", "the database host")
	// dbPort := flag.String("dbport", "5432", "the database port")
	// dbName := flag.String("dbname", "", "the database name")
	// dbUser := flag.String("dbuser", "", "the database user")
	// dbPassword := flag.String("dbpassword", "", "the database password")
	// dbSSL := flag.String("dbssl", "disable", "the database ssl settings(disable, prefer, require)")

	// flag.Parse()

	app.GoEnv = goEnv

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	validate := validator.New(validator.WithRequiredStructEnabled())
	app.Validate = validate

	// Connecto to DB
	log.Println("Connecting to dabase")
	connectionString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s", dbHost, dbPort, dbName, dbUser, dbPassword, dbSSL)
	db, err := driver.ConnectSQL(connectionString)
	if err != nil {
		log.Fatal("Cannot conect to database: Dying!", err)
	}
	log.Println("Connected to database")
	// 

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)

	return db, nil
}

func dir(envFile string) string {
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	for {
		goModPath := filepath.Join(currentDir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			break
		}

		parent := filepath.Dir(currentDir)
		if parent == currentDir {
			panic(fmt.Errorf("go.mod not found"))
		}
		currentDir = parent
	}

	return filepath.Join(currentDir, envFile)
}