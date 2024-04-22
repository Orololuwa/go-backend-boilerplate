package dbrepo

import (
	"context"
	"database/sql"
	"log"

	"github.com/Orololuwa/go-backend-boilerplate/src/repository"
)

type postgresDBRepo struct {
	DB *sql.DB
}
type testDBRepo struct {
	DB *sql.DB
}

func NewPostgresDBRepo(conn *sql.DB) repository.DatabaseRepo {
	return &postgresDBRepo{
		DB: conn,
	}
}
func NewTestingDBRepo() repository.DatabaseRepo {
	return &testDBRepo{
	}
}

func (m *postgresDBRepo) Transaction(ctx context.Context, operation func(context.Context, *sql.Tx) error) error {
    tx, err := m.DB.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
	
    defer func() error{
        if err != nil {
            tx.Rollback()
            return err
        }

        if err := tx.Commit(); err != nil {
            return err
        }

		log.Println("Transaction completed successfully")
        return nil
    }()

    if err := operation(ctx, tx); err != nil {
        return err
    }

    return nil
}