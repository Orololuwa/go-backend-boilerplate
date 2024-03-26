package repository

import "github.com/Orololuwa/go-backend-boilerplate/src/models"

type DatabaseRepo interface {
	GetHealth() bool

	InsertReservation(res models.Reservation) error
}