package dtos

type UserLoginBody struct {
	Email string `json:"email" validate:"required,email" faker:"email"`
}