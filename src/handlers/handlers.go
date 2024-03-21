package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Orololuwa/go-backend-boilerplate/src/config"
	"github.com/Orololuwa/go-backend-boilerplate/src/helpers"
)

type Repository struct {
	App *config.AppConfig
}

var Repo *Repository

// NewRepo function initializes the Repo
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
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