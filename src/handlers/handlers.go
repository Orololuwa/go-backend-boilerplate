package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Orololuwa/go-backend-boilerplate/src/config"
	"github.com/Orololuwa/go-backend-boilerplate/src/driver"
	"github.com/Orololuwa/go-backend-boilerplate/src/helpers"
	"github.com/Orololuwa/go-backend-boilerplate/src/models"
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

type ReservationBody struct {
	FirstName string `json:"firstName"`
	LastName string `json:"lastName"`
	Email string `json:"email"`
	Phone string `json:"phone"`
	StartDate string `json:"startDate"`
	EndDate string `json:"endDate"`
	RoomId string `json:"roomId"`
}

func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var body ReservationBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	sd := body.StartDate
	ed := body.EndDate

	layout := "2006-01-02"

	startDate, err := time.Parse(layout, sd)
	if err != nil {
		helpers.ServerError(w, err)
	}

	endDate, err := time.Parse(layout, ed)
	if err != nil {
		helpers.ServerError(w, err)
	}

	roomId, err := strconv.Atoi(body.RoomId)
	if err != nil {
		helpers.ServerError(w, err)
	}

	reservation := models.Reservation {
		FirstName: body.FirstName,
		LastName: body.LastName,
		Email: body.Email,
		Phone: body.Phone,
		StartDate: startDate,
		EndDate: endDate,
		RoomID: roomId,
	}

	log.Println(reservation)

	err = m.DB.InsertReservation(reservation)
	if err != nil {
		helpers.ServerError(w, err)
	}

	response := map[string]any{"message": "Request received successfully"}
    jsonResponse, err := json.Marshal(response)
    if err != nil {
        http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(jsonResponse)
}