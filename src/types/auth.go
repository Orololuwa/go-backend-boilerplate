package types

import "github.com/golang-jwt/jwt/v5"

type LoginSuccessResponse struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

type JWTClaims struct {
	Email string `json:"email"`
    jwt.RegisteredClaims
}