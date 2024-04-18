package dbrepo

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Orololuwa/go-backend-boilerplate/src/models"
)

func (m *postgresDBRepo) InsertReservation(ctx context.Context, tx *sql.Tx, res models.Reservation) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var newId int

	stmt := `insert into reservations (first_name, last_name, email, phone, start_date,
			 end_date, room_id, created_at, updated_at)
			 values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`

	var err error

	if tx != nil {
		err = tx.QueryRowContext(
			ctx, stmt, 
			res.FirstName,
			res.LastName,
			res.Email,
			res.Phone,
			res.StartDate,
			res.EndDate,
			res.RoomID,
			time.Now(),
			time.Now(),
		).Scan(&newId)
	}else{
		err = m.DB.QueryRowContext(
			ctx, stmt, 
			res.FirstName,
			res.LastName,
			res.Email,
			res.Phone,
			res.StartDate,
			res.EndDate,
			res.RoomID,
			time.Now(),
			time.Now(),
		).Scan(&newId)
	}

	
	
	if err != nil {
		return 0, err
	}

	return newId, nil
}

func (m *postgresDBRepo) InsertRoomRestriction(ctx context.Context, tx *sql.Tx, r models.RoomRestriction) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	stmt := `insert into room_restrictions (start_date, end_date, room_id, reservation_id,
			restriction_id, created_at, updated_at) values ($1, $2, $3, $4, $5, $6, $7 )`

	var err error
	if tx != nil {
		_, err = tx.ExecContext(
			ctx, stmt,
			r.StartDate,
			r.EndDate,
			r.RoomID,
			r.ReservationID,
			r.RestrictionID,
			time.Now(),
			time.Now(),
		)
	}else{
		_, err = m.DB.ExecContext(
			ctx, stmt,
			r.StartDate,
			r.EndDate,
			r.RoomID,
			r.ReservationID,
			r.RestrictionID,
			time.Now(),
			time.Now(),
		)
	}
	
	if err != nil {
		return err
	}

	return nil
}

// SearchAvailabilityForDatesByRoomId returns true if availability exists for a room_id and false if no availability exists
func (m *postgresDBRepo) SearchAvailabilityForDatesByRoomId(ctx context.Context, tx *sql.Tx, start, end time.Time, roomId int) (bool, error){
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var numRows int

	query := `
		select
			count(id)
		from
			room_restrictions
		where 
			room_id = $1
			and $2 < end_date and $3 > start_date
	`

	var row *sql.Row
	if tx != nil {
		row = tx.QueryRowContext(ctx, query, roomId, start, end)
	}else {
		row = m.DB.QueryRowContext(ctx, query, roomId, start, end)
	}
	err := row.Scan(&numRows)
	if err != nil {
		return false, err
	}

	if numRows == 0 {
		return true, nil
	}

	return false, nil
}

// SearchAvailabilityForAllRooms returns a slice of rooms for a given date range
func (m *postgresDBRepo) SearchAvailabilityForAllRooms(ctx context.Context, tx *sql.Tx, start, end time.Time) ([]models.Room, error){
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var rooms = make([]models.Room, 0)

	query := `
		select 
			r.id, r.room_name
		from
			rooms r
		where
			r.id not in 
		(select rr.room_id from room_restrictions rr where $1 < rr.end_date and $2 > rr.start_date)
	`

	var rows *sql.Rows
	var err error
	if tx != nil {
		rows, err = tx.QueryContext(ctx, query, start, end)
	}else{	
		rows, err = m.DB.QueryContext(ctx, query, start, end)
	}
	if err != nil {
		return rooms, err
	}

	for rows.Next(){
		var room models.Room
		err := rows.Scan(&room.ID, &room.RoomName)
		if err != nil {
			return rooms, err
		}
		rooms = append(rooms, room)
	}

	if err = rows.Err(); err != nil {
		return rooms, err
	}


	return rooms, nil
}

func (m *postgresDBRepo) GetRoomById(ctx context.Context, tx *sql.Tx, id int) (models.Room, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var room models.Room

	query := `
		select id, room_name, created_at, updated_at from rooms where id = $1
	`

	var row *sql.Row
	if tx != nil {
		row = tx.QueryRowContext(ctx, query, id)
	}else {
		row = m.DB.QueryRowContext(ctx, query, id)
	}
	err := row.Scan(&room.ID, &room.RoomName, &room.CreatedAt, &room.UpdatedAt)

	if err != nil {
		return room, err
	}

	return room, nil
}

func (m *postgresDBRepo) GetAllRooms(ctx context.Context, tx *sql.Tx, id int, room_name string, created_at string, updated_at string) ([]models.Room, error){
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var rooms = make([]models.Room, 0)

	query := `
		select 
			id, room_name, created_at, updated_at 
		from 
			rooms
		where
			1=1
	`
	args := []interface{}{}

	if id != 0 {
		query += fmt.Sprintf(" AND id = $%d", len(args)+1)
		args = append(args, id)
	}

	if room_name != "" {
		query += fmt.Sprintf(" AND room_name = $%d", len(args)+1)
		args = append(args, room_name)
	}

	if created_at != "" {
		query += fmt.Sprintf(" AND created_at = $%d", len(args)+1)
		args = append(args, created_at)
	}

	if updated_at != "" {
		query += fmt.Sprintf(" AND updated_at = $%d", len(args)+1)
		args = append(args, updated_at)
	}

	var rows *sql.Rows
	var err error
	if tx != nil {
		rows, err = tx.QueryContext(ctx, query, args...)
	}else{
		rows, err = m.DB.QueryContext(ctx, query, args...)
	}
	if err != nil {
        return rooms, err
    }

	for rows.Next(){
		var room models.Room
		err := rows.Scan(&room.ID, &room.RoomName, &room.CreatedAt, &room.UpdatedAt)
		if err != nil {
			return rooms, err
		}
		rooms = append(rooms, room)
	}

	if err = rows.Err(); err != nil {
		return rooms, err
	}

	return rooms, nil	
}