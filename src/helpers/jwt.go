package helpers

import (
	"errors"
	"os"
	"time"

	"github.com/Orololuwa/go-backend-boilerplate/src/types"
	"github.com/golang-jwt/jwt/v5"
)

func CreateJWTToken(email string) (string, error) {
	secretKey := []byte(os.Getenv("JWT_SECRET"))

	claims := types.JWTClaims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(10 * time.Minute)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
	  return "", err
	}
   
	return tokenString, nil
}

func VerifyJWTToken(tokenString string) (*jwt.Token, error) {
	secretKey := []byte(os.Getenv("JWT_SECRET"))

	token, err := jwt.ParseWithClaims(tokenString, &types.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
   
	if err != nil {
	   return nil, err
	}
   
	if !token.Valid {
	   return nil, errors.New("invalid token")
	}
   
	return token, nil
 }