package handlers

import (
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/Orololuwa/go-backend-boilerplate/src/config"
	"github.com/Orololuwa/go-backend-boilerplate/src/middleware"
	"github.com/go-chi/chi/v5"
	middlewareChi "github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
)

var testApp config.AppConfig
var mdTest *middleware.Middleware


func TestMain (m *testing.M){
	testApp.GoEnv = "test"

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	testApp.InfoLog = infoLog

	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	testApp.ErrorLog = errorLog

	validate := validator.New(validator.WithRequiredStructEnabled())
	testApp.Validate = validate

	repo := NewTestRepo(&testApp)
	NewHandlers(repo)

	mdTest = middleware.NewTest(&testApp)


	os.Exit(m.Run())
}


func getRoutes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middlewareChi.Logger)

	mux.Get("/health", Repo.Health)
	mux.Post("/reservation", Repo.PostReservation)
	mux.Post("/search-availability", Repo.SearchAvailability)
	mux.Post("/search-availability/{id}", Repo.SearchAvailabilityByRoomId)
	mux.Get("/room", Repo.GetAllRooms)
	mux.Get("/room/{id}", Repo.GetRoomById)

	return mux;
}