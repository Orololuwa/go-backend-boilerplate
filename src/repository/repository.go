package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/Orololuwa/go-backend-boilerplate/src/models"
)

type DatabaseRepo interface {
	Transaction(ctx context.Context, operation func(context.Context, *sql.Tx) error) error 
	InsertReservation(ctx context.Context, tx *sql.Tx, res models.Reservation) (int, error)
	InsertRoomRestriction(ctx context.Context, tx *sql.Tx, r models.RoomRestriction) error
	SearchAvailabilityForDatesByRoomId(ctx context.Context, tx *sql.Tx, start, end time.Time, roomId int) (bool, error)
	SearchAvailabilityForAllRooms(ctx context.Context, tx *sql.Tx, start, end time.Time) ([]models.Room, error)
	GetRoomById(ctx context.Context, tx *sql.Tx, id int) (models.Room, error)
	GetAllRooms(ctx context.Context, tx *sql.Tx, id int, room_name string, created_at string, updated_at string)([]models.Room, error)
}

type UserDBRepo interface {
	CreateAUser(ctx context.Context, tx *sql.Tx, user models.User) (int, error)
}