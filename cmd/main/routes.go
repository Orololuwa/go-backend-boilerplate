package main

import (
	"net/http"

	"github.com/Orololuwa/go-backend-boilerplate/src/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Logger)

	mux.Get("/health", handlers.Repo.Health)
	mux.Post("/reservation", handlers.Repo.PostReservation)
	mux.Post("/search-availability", handlers.Repo.PostAvailability)
	mux.Post("/search-availability/{id}", handlers.Repo.SearchAvailabilityByRoomId)
	mux.Get("/room", handlers.Repo.GetAllRooms)
	mux.Get("/room/{id}", handlers.Repo.GetRoomById)

	return mux;
}