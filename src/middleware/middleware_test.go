package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Orololuwa/go-backend-boilerplate/src/helpers"
	"github.com/go-faker/faker/v4"
)

type validationMiddleWareBody struct {
	Email string `json:"email" validate:"required,email" faker:"email"`
}

func middlewareHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func TestValidationMiddleware(t *testing.T){
	// test for missing body
	req := httptest.NewRequest("POST", "/route", nil)
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()

	reqBodyRef := &validationMiddleWareBody{}
	handlerChain := mdTest.ValidateReqBody(http.HandlerFunc(middlewareHandler), reqBodyRef)
	handlerChain.ServeHTTP(res, req)

	if res.Code != http.StatusBadRequest {
		t.Errorf("ValidateReqBody expected status code %d for missing request body, got %d", http.StatusBadRequest, res.Code)
	}

	// test for invalid email
	reqBody := validationMiddleWareBody{ Email: "johnDoe"}
	jsonData, err := json.Marshal(reqBody)
    if err != nil {
        t.Log("Error:", err)
        return
    }

	req = httptest.NewRequest("POST", "/route", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	res = httptest.NewRecorder()

	handlerChain = mdTest.ValidateReqBody(http.HandlerFunc(middlewareHandler), reqBodyRef)
	handlerChain.ServeHTTP(res, req)

	if res.Code != http.StatusBadRequest {
		t.Errorf("ValidateReqBody expected status code %d for invalid email, got %d", http.StatusBadRequest, res.Code)
	}

	// test for valid email
	err = faker.FakeData(&reqBody)
    if err != nil {
        t.Log(err)
    }
	jsonData, err = json.Marshal(reqBody)
    if err != nil {
        t.Log("Error:", err)
        return
    }

	req = httptest.NewRequest("POST", "/route", bytes.NewBuffer(jsonData))
	res = httptest.NewRecorder()

	handlerChain = mdTest.ValidateReqBody(http.HandlerFunc(middlewareHandler), reqBodyRef)
	handlerChain.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Errorf("ValidateReqBody expected status code %d, got %d", http.StatusOK, res.Code)
	}
}

func TestAuthorizationMiddleware(t *testing.T){
	// test for missing "Authorization" in request header
	req := httptest.NewRequest("POST", "/route", nil)
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()

	handlerChain := mdTest.Authorization(http.HandlerFunc(middlewareHandler))
	handlerChain.ServeHTTP(res, req)

	if res.Code != http.StatusUnauthorized {
		t.Errorf("Authorization expected status code %d for missing token in the header, got %d", http.StatusUnauthorized, res.Code)
	}

	// test for invalid or expired token
	// tokenString, err := helpers.CreateToken("johndoe@gmail.com")
	// if (err != nil){
	// 	t.Fatal("error creating test token")
	// }

	req = httptest.NewRequest("POST", "/route", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", "invalid tokenString"))
	res = httptest.NewRecorder()

	handlerChain = mdTest.Authorization(http.HandlerFunc(middlewareHandler))
	handlerChain.ServeHTTP(res, req)

	if res.Code != http.StatusUnauthorized {
		t.Errorf("Authorization expected status code %d for invalid token, got %d", http.StatusUnauthorized, res.Code)
	}

	// test for valid token
	tokenString, err := helpers.CreateJWTToken("johndoe@gmail.com")
	if (err != nil){
		t.Fatal("error creating test token")
	}

	req = httptest.NewRequest("POST", "/route", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokenString))
	res = httptest.NewRecorder()

	handlerChain = mdTest.Authorization(http.HandlerFunc(middlewareHandler))
	handlerChain.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Errorf("Authorization expected status code %d for invalid token, got %d", http.StatusOK, res.Code)
	}
}