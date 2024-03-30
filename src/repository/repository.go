package repository

import "github.com/Orololuwa/go-backend-boilerplate/src/models"

type DatabaseRepo interface {
	InsertReservation(res models.Reservation) (int, error)
	InsertRoomRestriction(r models.RoomRestriction) error
}