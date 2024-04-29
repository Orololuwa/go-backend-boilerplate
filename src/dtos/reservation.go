package dtos

type ReservationBody struct {
	FirstName string `json:"firstName" validate:"required" faker:"first_name"`
	LastName string `json:"lastName" validate:"required" faker:"last_name"`
	Email string `json:"email" validate:"required,email" faker:"email"`
	Phone string `json:"phone" validate:"required" faker:"phone_number"`
	StartDate string `json:"startDate" validate:"required" faker:"date"`
	EndDate string `json:"endDate" validate:"required" faker:"date"`
	RoomId int `json:"roomId" validate:"required" faker:"oneof: 15, 27, 61"`
}