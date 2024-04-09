package handlers

import (
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/Orololuwa/go-backend-boilerplate/src/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var testApp config.AppConfig

func TestMain (m *testing.M){
	testApp.GoEnv = "test"

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	testApp.InfoLog = infoLog

	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	testApp.ErrorLog = errorLog

	repo := NewTestRepo(&testApp)
	NewHandlers(repo)

	os.Exit(m.Run())
}


func getRoutes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Logger)

	mux.Get("/health", Repo.Health)
	mux.Post("/reservation", Repo.PostReservation)
	mux.Post("/search-availability", Repo.SearchAvailability)
	mux.Post("/search-availability/{id}", Repo.SearchAvailabilityByRoomId)
	mux.Get("/room", Repo.GetAllRooms)
	mux.Get("/room/{id}", Repo.GetRoomById)

	return mux;
}