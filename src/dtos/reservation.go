package dtos

type ReservationBody struct {
	FirstName string `json:"firstName" validate:"required"`
	LastName string `json:"lastName" validate:"required"`
	Email string `json:"email" validate:"required,email"`
	Phone string `json:"phone" validate:"required"`
	StartDate string `json:"startDate" validate:"required"`
	EndDate string `json:"endDate" validate:"required"`
	RoomId string `json:"roomId" validate:"required"`
}