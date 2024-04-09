package middleware

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type validationMiddleWareBody struct {
	Email string `json:"email" validate:"required,email"`
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
	handlerChain := ValidateMiddleware(http.HandlerFunc(middlewareHandler), reqBodyRef)

	handlerChain.ServeHTTP(res, req)

	if res.Code != http.StatusBadRequest {
		t.Errorf("ValidateMiddleware expected status code %d for missing request body, got %d", http.StatusBadRequest, res.Code)
	}

	// test for invalid email
	reqBody := validationMiddleWareBody{
		Email: "johndoe",
	}

	jsonData, err := json.Marshal(reqBody)
    if err != nil {
        t.Log("Error:", err)
        return
    }

	req = httptest.NewRequest("POST", "/route", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	res = httptest.NewRecorder()

	handlerChain = ValidateMiddleware(http.HandlerFunc(middlewareHandler), reqBodyRef)

	handlerChain.ServeHTTP(res, req)

	if res.Code != http.StatusBadRequest {
		t.Errorf("ValidateMiddleware expected status code %d for invalid email, got %d", http.StatusBadRequest, res.Code)
	}

	// test for valid email
	reqBody = validationMiddleWareBody{
		Email: "johnDoe@gmail.com",
	}

	jsonData, err = json.Marshal(reqBody)
    if err != nil {
        t.Log("Error:", err)
        return
    }

	req = httptest.NewRequest("POST", "/route", bytes.NewBuffer(jsonData))
	res = httptest.NewRecorder()

	handlerChain = ValidateMiddleware(http.HandlerFunc(middlewareHandler), reqBodyRef)

	handlerChain.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Errorf("ValidateMiddleware expected status code %d, got %d", http.StatusOK, res.Code)
	}
}