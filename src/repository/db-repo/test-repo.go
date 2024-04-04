package dbrepo

import (
	"errors"
	"time"

	"github.com/Orololuwa/go-backend-boilerplate/src/models"
)

func (m *testDBRepo) InsertReservation(res models.Reservation) (int, error) {
	// fail if roomId is 2
	if res.RoomID == 2 {
		return 0, errors.New("failed to insert reservation")
	}

	return 1, nil
}

func (m *testDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	// fail if i try to insert a room restriction for room id of 1000
	if r.RoomID == 1000 {
		return errors.New("failed to insert room restriction")
	}

 	return nil
}

// SearchAvailabilityForDatesByRoomId returns true if availability exists for a room_id and false if no availability exists
func (m *testDBRepo) SearchAvailabilityForDatesByRoomId(start, end time.Time, roomId int) (bool, error){
	// fail if roomId is 2
	if roomId == 2 {
		return false, errors.New("reservation for room not found")
	}

	return true, nil
}

// SearchAvailabilityForAllRooms returns a slice of rooms for a given date range
func (m *testDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error){
	var rooms = make([]models.Room, 0)

	// return an error when the year in the startDate is 1960
	if start.Year() < 1960 {
		return rooms, errors.New("error searching rooms")
	}

	return rooms, nil
}

func (m *testDBRepo) GetRoomById(id int) (models.Room, error) {
	var room models.Room

	if id == 1000 {
		return room, errors.New("error getting room")
	}

	return room, nil
}

func (m *testDBRepo) GetAllRooms(id int, room_name string, created_at string, updated_at string) ([]models.Room, error){
	var rooms = make([]models.Room, 0)

	if id == 2 {
		return rooms, errors.New("error getting rooms")
	}

	return rooms, nil	
}