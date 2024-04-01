package repository

import (
	"time"

	"github.com/Orololuwa/go-backend-boilerplate/src/models"
)

type DatabaseRepo interface {
	InsertReservation(res models.Reservation) (int, error)
	InsertRoomRestriction(r models.RoomRestriction) error
	SearchAvailabilityForDatesByRoomId(start, end time.Time, roomId int) (bool, error)
	SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error)
	GetRoomById(id int) (models.Room, error)
	GetAllRooms(id int, room_name string, created_at string, updated_at string)([]models.Room, error)
}