package dtos

type PostAvailabilityBody struct {
	StartDate string `json:"startDate" validate:"required"`
	EndDate string `json:"endDate" validate:"required"`
}