package main

import (
	"net/http"

	"github.com/Orololuwa/go-backend-boilerplate/src/config"
	"github.com/Orololuwa/go-backend-boilerplate/src/driver"
	"github.com/Orololuwa/go-backend-boilerplate/src/dtos"
	"github.com/Orololuwa/go-backend-boilerplate/src/handlers"
	middleware "github.com/Orololuwa/go-backend-boilerplate/src/middleware"
	"github.com/go-chi/chi/v5"
	middlewareChi "github.com/go-chi/chi/v5/middleware"
)

func routes(a *config.AppConfig, conn *driver.DB) http.Handler {
	// Initialize internal middlewares
	md := middleware.New(a, conn)	

	// 
	mux := chi.NewRouter()

	// middlewares
	mux.Use(middlewareChi.Logger)

	mux.Get("/health", handlers.Repo.Health)

	// reservations
	mux.Post("/reservation", handlers.Repo.PostReservation)

	// rooms
	mux.Post("/search-availability", md.ValidateReqBody(http.HandlerFunc(handlers.Repo.SearchAvailability), &dtos.PostAvailabilityBody{} ).ServeHTTP)
	mux.Post("/search-availability/{id}", md.ValidateReqBody(http.HandlerFunc(handlers.Repo.SearchAvailabilityByRoomId), &dtos.PostAvailabilityBody{}).ServeHTTP)
	mux.Get("/room", handlers.Repo.GetAllRooms)
	mux.Get("/room/{id}", handlers.Repo.GetRoomById)

	// auth
	mux.Post("/login", md.ValidateReqBody(http.HandlerFunc(handlers.Repo.LoginUser), &dtos.UserLoginBody{} ).ServeHTTP)

	// protected route
	mux.Get("/protected-route", md.Authorization(http.HandlerFunc(handlers.Repo.ProtectedRoute)).ServeHTTP)

	return mux;
}