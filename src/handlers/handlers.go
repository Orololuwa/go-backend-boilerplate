package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
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

// NewRepo function initializes the Repo
func NewTestRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
		DB: dbrepo.NewTestingDBRepo(),
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
		helpers.ClientError(w, err, http.StatusInternalServerError, "")
		return
	}

	sd := body.StartDate
	ed := body.EndDate

	layout := "2006-01-02"

	startDate, err := time.Parse(layout, sd)
	if err != nil {
		helpers.ClientError(w, err, http.StatusInternalServerError, "")
		return
	}

	endDate, err := time.Parse(layout, ed)
	if err != nil {
		helpers.ClientError(w, err, http.StatusInternalServerError, "")
		return
	}

	roomId, err := strconv.Atoi(body.RoomId)
	if err != nil {
		helpers.ClientError(w, err, http.StatusInternalServerError, "")
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
		helpers.ClientError(w, err, http.StatusInternalServerError, "")
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
		helpers.ClientError(w, err, http.StatusInternalServerError, "")
		return
	}

	helpers.ClientResponseWriter(w, nil, http.StatusCreated, "reservation booked successfully")
}

type PostAvailabilityBody struct {
	StartDate string `json:"startDate"`
	EndDate string `json:"endDate"`
}

func (m *Repository) SearchAvailability(w http.ResponseWriter, r *http.Request) {
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
		helpers.ClientError(w, err, http.StatusBadRequest, "")
		return
	}

	endDate, err := time.Parse(layout, end)
	if err != nil {
		helpers.ClientError(w, err, http.StatusBadRequest, "")
		return
	}

	rooms, err := m.DB.SearchAvailabilityForAllRooms(startDate, endDate)
	if (err != nil){
		helpers.ClientError(w, err, http.StatusNotFound, "")
		return
	}

	helpers.ClientResponseWriter(w, rooms, http.StatusFound, "rooms retrieved successfully")
}

func (m *Repository) SearchAvailabilityByRoomId(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var id int
	var err error
	if m.App.GoEnv == "test" {
		exploded := strings.Split(r.RequestURI, "/")
		id, err = strconv.Atoi(exploded[2])
		if err != nil {
			helpers.ClientError(w, err, http.StatusInternalServerError, "missing URL param")
		return
		}
	}else{		
		id, err = strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			helpers.ClientError(w, err, http.StatusInternalServerError, "")
			return
		}
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
		helpers.ClientError(w, err, http.StatusNotFound, "")
		return
	}

	helpers.ClientResponseWriter(w, isRoomAvailable, http.StatusFound, "room retrieved successfully")

}

func (m *Repository) GetRoomById(w http.ResponseWriter, r *http.Request){
	var id int
	var err error
	if m.App.GoEnv == "test" {
		exploded := strings.Split(r.RequestURI, "/")
		id, err = strconv.Atoi(exploded[2])
		if err != nil {
			helpers.ClientError(w, err, http.StatusInternalServerError, "missing URL param")
		return
		}
	}else{		
		id, err = strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			helpers.ClientError(w, err, http.StatusInternalServerError, "")
			return
		}
	}

	room, err := m.DB.GetRoomById(id)
	if err != nil {
		helpers.ClientError(w, err, http.StatusNotFound, "room not found")
		return
	}

	helpers.ClientResponseWriter(w, room, http.StatusOK, "room retrieved successfully")
}

func (m *Repository) GetAllRooms(w http.ResponseWriter, r *http.Request){
	var room_name string
	var id int

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

	helpers.ClientResponseWriter(w, rooms, http.StatusOK, "rooms retrieved successfully")
}

