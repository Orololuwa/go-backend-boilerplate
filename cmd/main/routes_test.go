package main

import (
	"database/sql"
	"fmt"
	"log"
	"testing"

	"github.com/Orololuwa/go-backend-boilerplate/src/driver"
	"github.com/go-chi/chi/v5"
)

func TestRoutes(t *testing.T){
	sql := createTestDBInstance()
	conn := driver.DB{SQL: sql}
	mux := routes(&testApp, &conn)

	switch v := mux.(type) {
	case *chi.Mux:
		// do nothing
	default:
		t.Errorf(fmt.Sprintf("type is not *chi.Mux, but is %T", v))
	}
}

func createTestDBInstance() *sql.DB {
    db, err := sql.Open("pgx", ":memory:")
    if err != nil {
        log.Fatalf("Error opening database connection: %v", err)
    }
	
    return db
}