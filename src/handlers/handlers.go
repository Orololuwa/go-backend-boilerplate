package handlers

import (
	"encoding/json"
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
	Message string `json:"message"`
	Data interface{} `json:"data"`
}

func (m *Repository) Health(w http.ResponseWriter, r *http.Request){
	resp := jsonResponse{
		Message: "App Healthy",
		Data: nil,
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
		return
	}

	endDate, err := time.Parse(layout, ed)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	roomId, err := strconv.Atoi(body.RoomId)
	if err != nil {
		helpers.ServerError(w, err)
		return
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


	newReservationId, err := m.DB.InsertReservation(reservation)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	restriction := models.RoomRestriction{
		StartDate: startDate,
		EndDate: endDate,
		RoomID: roomId,
		ReservationID: newReservationId,
		RestrictionID: 1,
	}

	err = m.DB.InsertRoomRestriction(restriction)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	response := map[string]any{"message": "Request received successfully", "data": nil}
    jsonResponse, err := json.Marshal(response)
    if err != nil {
        http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(jsonResponse)
}

type PostAvailabilityBody struct {
	StartDate string `json:"startDate"`
	EndDate string `json:"endDate"`
}

func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var body PostAvailabilityBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	start := body.StartDate
	end := body.EndDate

	layout := "2006-01-02"

	startDate, err := time.Parse(layout, start)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	endDate, err := time.Parse(layout, end)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	rooms, err := m.DB.SearchAvailabilityForAllRooms(startDate, endDate)
	if (err != nil){
		helpers.ServerError(w, err)
		return
	}

	response := map[string]interface{}{"message": "rooms retrieved successfully", "data": rooms}
    jsonResponse, err := json.Marshal(response)
    if err != nil {
        http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusFound)
    w.Write(jsonResponse)
}