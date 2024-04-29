package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Orololuwa/go-backend-boilerplate/src/config"
	"github.com/Orololuwa/go-backend-boilerplate/src/driver"
	"github.com/Orololuwa/go-backend-boilerplate/src/dtos"
	"github.com/Orololuwa/go-backend-boilerplate/src/helpers"
	"github.com/Orololuwa/go-backend-boilerplate/src/models"
	"github.com/Orololuwa/go-backend-boilerplate/src/repository"
	dbrepo "github.com/Orololuwa/go-backend-boilerplate/src/repository/db-repo"
	"github.com/Orololuwa/go-backend-boilerplate/src/types"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type Repository struct {
	App *config.AppConfig
	DB repository.DatabaseRepo
	User repository.UserDBRepo
}

var Repo *Repository

// NewRepo function initializes the Repo
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB: dbrepo.NewPostgresDBRepo(db.SQL),
		User: dbrepo.NewUserDBRepo(db.SQL),		
	}
}

// NewRepo function initializes the Repo
func NewTestRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
		DB: dbrepo.NewTestingDBRepo(),
		User: dbrepo.NewUserTestingDBRepo(),
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

func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var body dtos.ReservationBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		helpers.ClientError(w, err, http.StatusInternalServerError, "")
		return
	}


	// validate the request body
	err = m.App.Validate.Struct(body)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		log.Println(err)
		helpers.ClientError(w, err, http.StatusBadRequest, errors.Error())
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

	// roomId, err := strconv.Atoi(body.RoomId)
	// if err != nil {
	// 	helpers.ClientError(w, err, http.StatusInternalServerError, "")
	// 	return
	// }

	reservation := models.Reservation {
		FirstName: body.FirstName,
		LastName: body.LastName,
		Email: body.Email,
		Phone: body.Phone,
		StartDate: startDate,
		EndDate: endDate,
		RoomID: body.RoomId,
	}

	ctx := context.Background()

	err = m.DB.Transaction(ctx, func(ctx context.Context, tx *sql.Tx) error {
		newReservationId, err := m.DB.InsertReservation(ctx, tx, reservation)
		if err != nil {
            return err
        }
	
		restriction := models.RoomRestriction{
			StartDate: startDate,
			EndDate: endDate,
			RoomID: body.RoomId,
			ReservationID: newReservationId,
			RestrictionID: 1,
		}
	
		err = m.DB.InsertRoomRestriction(ctx, tx, restriction)
		if err != nil {
            return err
        }

		return nil
	})


	if err != nil {
		helpers.ClientError(w, err, http.StatusInternalServerError, "")
		return
	}



	helpers.ClientResponseWriter(w, nil, http.StatusCreated, "reservation booked successfully")
}

func (m *Repository) SearchAvailability(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var body dtos.PostAvailabilityBody
	requestBody, ok := r.Context().Value("validatedRequestBody").(*dtos.PostAvailabilityBody)
    if !ok || requestBody == nil {
		helpers.ClientError(w, errors.New("failed to retrieve request body"), http.StatusBadRequest, "")
        return
    }
	body = *requestBody

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

	ctx := context.Background()

	rooms, err := m.DB.SearchAvailabilityForAllRooms(ctx, nil, startDate, endDate)
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
	
	exploded := strings.Split(r.RequestURI, "/")
	id, err = strconv.Atoi(exploded[2])
	if err != nil {
		helpers.ClientError(w, err, http.StatusInternalServerError, "missing URL param")
		return
	}

	var body dtos.PostAvailabilityBody
	requestBody, ok := r.Context().Value("validatedRequestBody").(*dtos.PostAvailabilityBody)
    if !ok || requestBody == nil {
		helpers.ClientError(w, errors.New("failed to retrieve request body"), http.StatusBadRequest, "")
        return
    }
	body = *requestBody

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

	ctx := context.Background()


	isRoomAvailable, err := m.DB.SearchAvailabilityForDatesByRoomId(ctx, nil, startDate, endDate, id)
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

	room, err := m.DB.GetRoomById(context.Background(), nil, id)
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

	rooms, err := m.DB.GetAllRooms(context.Background(), nil, id, room_name, "", "")
	if err != nil {
		helpers.ClientError(w, err, http.StatusNotFound, "")
		return
	}

	helpers.ClientResponseWriter(w, rooms, http.StatusOK, "rooms retrieved successfully")
}


func (m *Repository) LoginUser(w http.ResponseWriter, r *http.Request){
	var body dtos.UserLoginBody
	requestBody, ok := r.Context().Value("validatedRequestBody").(*dtos.UserLoginBody)
    if !ok || requestBody == nil {
		helpers.ClientError(w, errors.New("failed to retrieve request body"), http.StatusBadRequest, "")
        return
    }
	body = *requestBody

	tokenString, err := helpers.CreateJWTToken(body.Email)

	if err != nil {
		helpers.ClientError(w, err, http.StatusInternalServerError, "")
	}

	data := types.LoginSuccessResponse{Email: body.Email, Token: tokenString}

	helpers.ClientResponseWriter(w, data, http.StatusOK, "logged in successfully")
}

func (m *Repository) ProtectedRoute(w http.ResponseWriter, r *http.Request){
	helpers.ClientResponseWriter(w, nil, http.StatusOK, "welcome to the protected route")
}
