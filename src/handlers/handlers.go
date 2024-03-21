package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Orololuwa/go-backend-boilerplate/src/config"
	"github.com/Orololuwa/go-backend-boilerplate/src/driver"
	"github.com/Orololuwa/go-backend-boilerplate/src/helpers"
	"github.com/Orololuwa/go-backend-boilerplate/src/repository"
	dbrepo "github.com/Orololuwa/go-backend-boilerplate/src/repository/db-repo"
)

type Repository struct {
	App *config.AppConfig
	DB repository.DatabaseRepo
}

var Repo *Repository

// NewRepo function initializes the Repo
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB: dbrepo.NewPostgresDBRepo(db.SQL),
	}
}

func NewHandlers(r *Repository){
	Repo = r;
}

type jsonResponse struct {
	Status int `json:"status"`
	Message string `json:"message"`
	Data string `json:"data"`
}

func (m *Repository) Health(w http.ResponseWriter, r *http.Request){
	isHealthOk := m.DB.GetHealth()
	log.Println(isHealthOk)

	resp := jsonResponse{
		Status: http.StatusOK,
		Message: "App Healthy",
		Data: "OK!",
	}

	out, err := json.MarshalIndent(resp, "", "     ")
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}