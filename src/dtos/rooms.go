package dtos

type PostAvailabilityBody struct {
	StartDate string `json:"startDate" validate:"required" faker:"date"`
	EndDate string `json:"endDate" validate:"required" faker:"date"`
}