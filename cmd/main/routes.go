package main

import (
	"net/http"

	"github.com/Orololuwa/go-backend-boilerplate/src/dtos"
	"github.com/Orololuwa/go-backend-boilerplate/src/handlers"
	middlewareInternal "github.com/Orololuwa/go-backend-boilerplate/src/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func routes() http.Handler {
	// 
	mux := chi.NewRouter()

	// middlewares
	mux.Use(middleware.Logger)

	mux.Get("/health", handlers.Repo.Health)

	// reservations
	mux.Post("/reservation", handlers.Repo.PostReservation)

	// rooms
	mux.Post("/search-availability", middlewareInternal.ValidateReqBody(http.HandlerFunc(handlers.Repo.SearchAvailability), &dtos.PostAvailabilityBody{} ).ServeHTTP)
	mux.Post("/search-availability/{id}", middlewareInternal.ValidateReqBody(http.HandlerFunc(handlers.Repo.SearchAvailabilityByRoomId), &dtos.PostAvailabilityBody{} ).ServeHTTP)
	mux.Get("/room", handlers.Repo.GetAllRooms)
	mux.Get("/room/{id}", handlers.Repo.GetRoomById)

	return mux;
}