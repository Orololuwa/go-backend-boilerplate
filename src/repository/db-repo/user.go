package dbrepo

import (
	"context"
	"database/sql"
	"time"

	"github.com/Orololuwa/go-backend-boilerplate/src/models"
	"github.com/Orololuwa/go-backend-boilerplate/src/repository"
)

type user struct {
	DB *sql.DB
}
func NewUserDBRepo(conn *sql.DB) repository.UserDBRepo {
	return &user{
		DB: conn,
	}
}

type testUserDBRepo struct {
	DB *sql.DB
}
func NewUserTestingDBRepo() repository.UserDBRepo {
	return &testUserDBRepo{
	}
}

func (m *user) CreateAUser(user models.User) (int, error){
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var newId int

	query := `
			INSERT into users 
				(first_name, last_name, email, password, created_at, updated_at)
			values 
				($1, $2, $3, $4, $5, $6)
			returning id`
	
	err := m.DB.QueryRowContext(ctx, query, 
		user.FirstName, 
		user.LastName, 
		user.Email, 
		user.Password,
		time.Now(),
		time.Now(),
	).Scan(&newId)

	if err != nil {
		return 0, err
	}

	return newId, nil
}

func (m *user) GetAUser(id int) (models.User, error){
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user models.User

	query := `
			SELECT (id, first_name, last_name, email, password, created_at, updated_at)
			from users
			WHERE
			id=$1
	`

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return user, err
	}

	return user, nil
}

func (m *user) GetAllUser() ([]models.User, error){
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var users = make([]models.User, 0)

	query := `
		SELECT (id, first_name, last_name, email, password, created_at, updated_at)
		from users
	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return users, err
	}

	for rows.Next(){
		var user models.User
		err := rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&user.Password,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return users, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return users, err
	}

	return users, nil
}

func (m *user) UpdateAUsersName(id int, firstName, lastName string)(error){
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		UPDATE 
			users set (first_name, last_name) = ($1, $2)
		WHERE
			id = $3
	`

	_, err := m.DB.ExecContext(ctx, query, firstName, lastName, id)
	if err != nil{
		return  err
	}

	return nil
}

func (m *user) DeleteUserByID(id int) error {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    query := "DELETE FROM users WHERE id = $1"

    _, err := m.DB.ExecContext(ctx, query, id)
    if err != nil {
        return err
    }

    return nil
}