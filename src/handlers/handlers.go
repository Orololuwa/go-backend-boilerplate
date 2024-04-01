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
	"github.com/go-chi/chi/v5"
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
		helpers.ClientError(w, err, http.StatusBadRequest, "")
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

func (m *Repository) SearchAvailabilityByRoomId(w http.ResponseWriter, r *http.Request){
	// var start_date, end_date, id string
	// start_date = r.URL.Query().Get("start_date")
	// end_date = r.URL.Query().Get("end_date")
	// id = r.URL.Query().Get("id")

	// log.Println(start_date, end_date, id)

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	var body PostAvailabilityBody
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		helpers.ClientError(w, err, http.StatusBadRequest, "")
		return
	}

	start := body.StartDate
	end := body.EndDate

	layout := "2006-01-02"

	startDate, err := time.Parse(layout, start)
	if err != nil {
		helpers.ClientError(w, err, http.StatusBadRequest, "")
		return
	}

	endDate, err := time.Parse(layout, end)
	if err != nil {
		helpers.ClientError(w, err, http.StatusBadRequest, "")
		return
	}

	isRoomAvailable, err := m.DB.SearchAvailabilityForDatesByRoomId(startDate, endDate, id)
	if err != nil {
		helpers.ClientError(w, err, http.StatusBadRequest, "")
		return
	}

	response := map[string]interface{}{"message": "room retrieved successfully", "data": isRoomAvailable}
    jsonResponse, err := json.Marshal(response)
    if err != nil {
        http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusFound)
    w.Write(jsonResponse)
}

func (m *Repository) GetRoomById(w http.ResponseWriter, r *http.Request){
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	room, err := m.DB.GetRoomById(id)
	if err != nil {
		helpers.ClientError(w, err, http.StatusNotFound, "room not found")
		return
	}

	response := map[string]interface{}{"message": "room retrieved successfully", "data": room}
    jsonResponse, err := json.Marshal(response)
    if err != nil {
        http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusFound)
    w.Write(jsonResponse)
}

func (m *Repository) GetAllRooms(w http.ResponseWriter, r *http.Request){
	var room_name string
	var id int
	// start_date = r.URL.Query().Get("start_date")

	room_name = r.URL.Query().Get("room_name")

	if r.URL.Query().Has("id"){
		paramId, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			helpers.ClientError(w, err, http.StatusNotFound, "")
			return
		}
		id = paramId
	}


	rooms, err := m.DB.GetAllRooms(id, room_name, "", "")
	if err != nil {
		helpers.ClientError(w, err, http.StatusNotFound, "")
		return
	}

	response := map[string]interface{}{"message": "room retrieved successfully", "data": rooms}
    jsonResponse, err := json.Marshal(response)
    if err != nil {
        http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusFound)
    w.Write(jsonResponse)
}

