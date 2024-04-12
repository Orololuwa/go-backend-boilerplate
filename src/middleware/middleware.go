package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Orololuwa/go-backend-boilerplate/src/helpers"
	"github.com/Orololuwa/go-backend-boilerplate/src/types"
	"github.com/go-playground/validator/v10"
)

func ValidateReqBody(next http.Handler, requestBodyStruct interface{}) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

        decoder := json.NewDecoder(r.Body)
        if err := decoder.Decode(requestBodyStruct); err != nil {
			helpers.ClientError(w, err, http.StatusBadRequest, "failed to decode body")
            return
        }

		defer r.Body.Close()

        validate := validator.New()
        if err := validate.Struct(requestBodyStruct); err != nil {
            errors := err.(validator.ValidationErrors)
			helpers.ClientError(w, err, http.StatusBadRequest, errors.Error())
            return
        }

		ctx := context.WithValue(r.Context(), "validatedRequestBody", requestBodyStruct)
		r = r.WithContext(ctx)

        next.ServeHTTP(w, r)
    })
}

func Authorization(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        tokenString := r.Header.Get("Authorization")
        if tokenString == "" {
            helpers.ClientError(w, errors.New("missing token"), http.StatusUnauthorized, "")
            return
        }
        tokenString = tokenString[len("Bearer "):]

        token, err := helpers.VerifyJWTToken(tokenString)
        if err != nil {
            helpers.ClientError(w, errors.New("invalid or expired token"), http.StatusUnauthorized, "")
            return
        }

        claims, ok := token.Claims.(*types.JWTClaims)
        if ok {
            // get the user's data and verify
            fmt.Println(claims.Email)
        }else{
            helpers.ClientError(w, errors.New("unknown claims type, cannot proceed"), http.StatusInternalServerError, "")
            return
        }

        next.ServeHTTP(w, r)
    })
}