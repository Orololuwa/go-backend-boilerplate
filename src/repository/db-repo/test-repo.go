package dbrepo

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Orololuwa/go-backend-boilerplate/src/models"
)

// Transactions
func (m *testDBRepo) Transaction(ctx context.Context, operation func(context.Context, *sql.Tx) error) error {
	if err := operation(ctx, nil); err != nil {
        return err
    }

    return nil
}

// Reservations
func (m *testDBRepo) InsertReservation(ctx context.Context, tx *sql.Tx, res models.Reservation) (int, error) {
	// fail if roomId is 2
	if res.RoomID == 2 {
		return 0, errors.New("failed to insert reservation")
	}

	return 1, nil
}

// Room restrictions
func (m *testDBRepo) InsertRoomRestriction(ctx context.Context, tx *sql.Tx, r models.RoomRestriction) error {
	// fail if i try to insert a room restriction for room id of 1000
	if r.RoomID == 1000 {
		return errors.New("failed to insert room restriction")
	}

 	return nil
}

// Rooms
// SearchAvailabilityForAllRooms returns a slice of rooms for a given date range
func (m *testDBRepo) SearchAvailabilityForAllRooms(ctx context.Context, tx *sql.Tx, start, end time.Time) ([]models.Room, error){
	var rooms = make([]models.Room, 0)

	// return an error when the year in the startDate is 1960
	if start.Year() < 1960 {
		return rooms, errors.New("error searching rooms")
	}

	return rooms, nil
}

// SearchAvailabilityForDatesByRoomId returns true if availability exists for a room_id and false if no availability exists
func (m *testDBRepo) SearchAvailabilityForDatesByRoomId(ctx context.Context, tx *sql.Tx, start, end time.Time, roomId int) (bool, error){
	// simulate failure for roomId 2
	if roomId == 2 {
		return false, errors.New("reservation for room not found")
	}

	return true, nil
}

func (m *testDBRepo) GetAllRooms(ctx context.Context, tx *sql.Tx, id int, room_name string, created_at string, updated_at string) ([]models.Room, error){
	var rooms = make([]models.Room, 0)

	// simulate failure for roomId 2
	if id == 2 {
		return rooms, errors.New("error getting rooms")
	}

	return rooms, nil	
}

func (m *testDBRepo) GetRoomById(ctx context.Context, tx *sql.Tx, id int) (models.Room, error) {
	var room models.Room

	if id == 1000 {
		return room, errors.New("error getting room")
	}

	return room, nil
}

// User
func (m *testUserDBRepo) CreateAUser(ctx context.Context, tx *sql.Tx, user models.User) (int, error){
	var newId int


	return newId, nil
}