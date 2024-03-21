package repository

type DatabaseRepo interface {
	GetHealth() bool
}