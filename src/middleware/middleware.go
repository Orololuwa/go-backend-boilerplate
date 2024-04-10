package middleware

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Orololuwa/go-backend-boilerplate/src/helpers"
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