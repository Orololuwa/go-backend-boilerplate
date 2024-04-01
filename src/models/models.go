package models

import "time"

// Users is the user's model
type User struct {
	ID int
	FirstName string
	LastName  string
	Email     string
	Password string
	AccessLevel int
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Room struct {
	ID int `json:"id"`
	RoomName string `json:"roomName"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

}

type Restriction struct {
	ID int
	RestrictionName string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Reservation struct {
	ID int
	FirstName string
	LastName  string
	Email     string
	Phone string
	StartDate time.Time
	EndDate time.Time
	RoomID int
	CreatedAt time.Time
	UpdatedAt time.Time
	Room Room
}

type RoomRestriction struct {
	ID int
	StartDate time.Time
	EndDate time.Time
	RoomID int
	ReservationID int
	RestrictionID int
	CreatedAt time.Time
	UpdatedAt time.Time
	Room Room
	Reservation Reservation
	Restriction Restriction
}
