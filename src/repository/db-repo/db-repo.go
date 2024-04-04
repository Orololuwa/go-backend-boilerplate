package dbrepo

import (
	"database/sql"

	"github.com/Orololuwa/go-backend-boilerplate/src/repository"
)

type postgresDBRepo struct {
	DB *sql.DB
}

type testDBRepo struct {
	DB *sql.DB
}

func NewPostgresDBRepo (conn *sql.DB) repository.DatabaseRepo {
	return &postgresDBRepo{
		DB: conn,
	}
}

func NewTestingDBRepo () repository.DatabaseRepo {
	return &testDBRepo{
	}
}