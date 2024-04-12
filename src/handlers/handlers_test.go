package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/Orololuwa/go-backend-boilerplate/src/dtos"
	"github.com/Orololuwa/go-backend-boilerplate/src/middleware"
)

type postData struct {
	key string
	value string
}

var theTests = []struct {
	name string
	url string
	method string
	params []postData
	expectedStatusCode int
}{
	{"health", "/health", "GET", []postData{}, http.StatusOK},
	// {"post reservation", "/reservation", "POST", []postData{}, http.StatusOK},
	// {"search availability", "/search-availability", "POST", []postData{}, http.StatusOK},
	// {"search availability by room_id", "/search-availability/{id}", "POST", []postData{}, http.StatusOK},
	// {"get rooms", "/room", "GET", []postData{}, http.StatusOK},
	// {"get room by id", "/room/{id}", "GET", []postData{}, http.StatusOK},
}

func TestHandler(t *testing.T){
	routes := getRoutes()
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, e := range theTests {
		if e.method == "GET" {
			resp, err := ts.Client().Get(ts.URL + e.url)
			if (err != nil){
				t.Log(err)
				t.Fatal(err)
			}

			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		}
		// else{
		// 	// do something else
		// }
	}
}

func TestRepository_PostReservation(t *testing.T){
	// test if I try to call a method other than POST
	reqBody := []byte(`
	{
		"firstName": "John",
		"lastName": "Doe",
		"email": "johndoe@go.lang",
		"phone": "+234-000-000-0000",
		"startDate": "2024-03-30",
		"endDate": "2024-04-30",
		"roomId": "1"
	}
	`)
	req, _ := http.NewRequest("PUT", "/reservation", bytes.NewBuffer([]byte(reqBody)))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("PostReservation handler returned wrong response code for wrong http method: got %d, wanted %d", rr.Code, http.StatusMethodNotAllowed)
	}

	// test for the right request body
	req, _ = http.NewRequest("POST", "/reservation", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("PostReservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusCreated)
	}


	// test for missing body
	req, _ = http.NewRequest("POST", "/reservation", bytes.NewBuffer([]byte(``)))
	req.Header.Set("Content-Type", "application/json")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("PostReservation handler returned wrong response code for missing body: got %d, wanted %d", rr.Code, http.StatusInternalServerError)
	}

	// test for validator for invalid email
	reqBody = []byte(`
	{
		"firstName": "John",
		"lastName": "Doe",
		"email": "johndoe",
		"phone": "+234-000-000-0000",
		"startDate": "invalid",
		"endDate": "2024-04-30",
		"roomId": "1"
	}
	`)

	req, _ = http.NewRequest("POST", "/reservation", bytes.NewBuffer([]byte(reqBody)))
	req.Header.Set("Content-Type", "application/json")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("PostReservation handler returned wrong response code for invalid email on validator: got %d, wanted %d", rr.Code, http.StatusBadRequest)
	}

	// test for invalid start date
	reqBody = []byte(`
	{
		"firstName": "John",
		"lastName": "Doe",
		"email": "johndoe@go.lang",
		"phone": "+234-000-000-0000",
		"startDate": "invalid",
		"endDate": "2024-04-30",
		"roomId": "1"
	}
	`)

	req, _ = http.NewRequest("POST", "/reservation", bytes.NewBuffer([]byte(reqBody)))
	req.Header.Set("Content-Type", "application/json")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("PostReservation handler returned wrong response code for invalid startDate: got %d, wanted %d", rr.Code, http.StatusInternalServerError)
	}

	// test for invalid end date
	reqBody = []byte(`
	{
		"firstName": "John",
		"lastName": "Doe",
		"email": "johndoe@go.lang",
		"phone": "+234-000-000-0000",
		"startDate": "2024-04-30",
		"endDate": "invalid",
		"roomId": "1"
	}
	`)

	req, _ = http.NewRequest("POST", "/reservation", bytes.NewBuffer([]byte(reqBody)))
	req.Header.Set("Content-Type", "application/json")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("PostReservation handler returned wrong response code for invalid endDate: got %d, wanted %d", rr.Code, http.StatusInternalServerError)
	}

	// test for invalid roomId
	reqBody = []byte(`
	{
		"firstName": "John",
		"lastName": "Doe",
		"email": "johndoe@go.lang",
		"phone": "+234-000-000-0000",
		"startDate": "2024-04-30",
		"endDate": "2024-04-30",
		"roomId": "one"
	}
	`)

	req, _ = http.NewRequest("POST", "/reservation", bytes.NewBuffer([]byte(reqBody)))
	req.Header.Set("Content-Type", "application/json")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("PostReservation handler returned wrong response code for invalid roomId: got %d, wanted %d", rr.Code, http.StatusInternalServerError)
	}

	// test for failure to insert the reservation
	reqBody = []byte(`
	{
		"firstName": "John",
		"lastName": "Doe",
		"email": "johndoe@go.lang",
		"phone": "+234-000-000-0000",
		"startDate": "2024-04-30",
		"endDate": "2024-04-30",
		"roomId": "2"
	}
	`)

	req, _ = http.NewRequest("POST", "/reservation", bytes.NewBuffer([]byte(reqBody)))
	req.Header.Set("Content-Type", "application/json")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("PostReservation handler returned wrong response code for failed DB insert: got %d, wanted %d", rr.Code, http.StatusInternalServerError)
	}

	// test for failure to insert room restriction
	reqBody = []byte(`
	{
		"firstName": "John",
		"lastName": "Doe",
		"email": "johndoe@go.lang",
		"phone": "+234-000-000-0000",
		"startDate": "2024-04-30",
		"endDate": "2024-04-30",
		"roomId": "1000"
	}
	`)

	req, _ = http.NewRequest("POST", "/reservation", bytes.NewBuffer([]byte(reqBody)))
	req.Header.Set("Content-Type", "application/json")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("PostReservation handler returned wrong response code for failed DB insert: got %d, wanted %d", rr.Code, http.StatusInternalServerError)
	}
}

func TestRepository_SearchAvailability(t *testing.T){
	// Test to make sure that a post handler is being called
	reqBody := dtos.PostAvailabilityBody{
		StartDate: "2024-05-30",
		EndDate: "2024-06-06",
	}
	
	jsonData, err := json.Marshal(reqBody)
    if err != nil {
        t.Log("Error:", err)
        return
    }

	req, _ := http.NewRequest("PUT", "/search-availability", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	res := httptest.NewRecorder()

	reqBodyRef := &dtos.PostAvailabilityBody{}
	handlerChain := middleware.ValidateReqBody(http.HandlerFunc(Repo.SearchAvailability), reqBodyRef)

	handlerChain.ServeHTTP(res, req)

	if res.Code != http.StatusMethodNotAllowed {
		t.Errorf("SearchAvailability handler returned wrong response code for wrong http method: got %d, wanted %d", res.Code, http.StatusMethodNotAllowed)
	}

	// test for the right request body
	reqBody = dtos.PostAvailabilityBody{
		StartDate: "2024-05-30",
		EndDate: "2024-06-06",
	}
	
	jsonData, err = json.Marshal(reqBody)
    if err != nil {
        t.Log("Error:", err)
        return
    }

	req, _ = http.NewRequest("POST", "/search-availability", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	res = httptest.NewRecorder()

	reqBodyRef = &dtos.PostAvailabilityBody{}
	handlerChain = middleware.ValidateReqBody(http.HandlerFunc(Repo.SearchAvailability), reqBodyRef)

	handlerChain.ServeHTTP(res, req)

	if res.Code != http.StatusFound {
		t.Errorf("SearchAvailability handler returned wrong response code: got %d, wanted %d", res.Code, http.StatusFound)
	}

	// test for missing request body
	req, _ = http.NewRequest("POST", "/search-availability", bytes.NewBuffer([]byte(``)))
	req.Header.Set("Content-Type", "application/json")

	res = httptest.NewRecorder()

	reqBodyRef = &dtos.PostAvailabilityBody{}
	handlerChain = middleware.ValidateReqBody(http.HandlerFunc(Repo.SearchAvailability), reqBodyRef)

	handlerChain.ServeHTTP(res, req)

	if res.Code != http.StatusBadRequest {
		t.Errorf("SearchAvailability handler returned wrong response code for missing request body: got %d, wanted %d", res.Code, http.StatusBadRequest)
	}

	// test for missing request body in context
	req, _ = http.NewRequest("POST", "/search-availability", bytes.NewBuffer([]byte(``)))
	req.Header.Set("Content-Type", "application/json")

	res = httptest.NewRecorder()

	handlerChain = http.HandlerFunc(Repo.SearchAvailability)

	handlerChain.ServeHTTP(res, req)

	if res.Code != http.StatusBadRequest {
		t.Errorf("SearchAvailability handler returned wrong response code for missing request body: got %d, wanted %d", res.Code, http.StatusBadRequest)
	}

	// test for invalid startDate
	reqBody = dtos.PostAvailabilityBody{
		StartDate: "invalid",
		EndDate: "2024-06-06",
	}
	
	jsonData, err = json.Marshal(reqBody)

	if err != nil {
		t.Log("Error:", err)
		return
	}
	req, _ = http.NewRequest("POST", "/search-availability", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	res = httptest.NewRecorder()

	reqBodyRef = &dtos.PostAvailabilityBody{}
	handlerChain = middleware.ValidateReqBody(http.HandlerFunc(Repo.SearchAvailability), reqBodyRef)

	handlerChain.ServeHTTP(res, req)

	if res.Code != http.StatusBadRequest {
		t.Errorf("SearchAvailability handler returned wrong response code for invalid startDate: got %d, wanted %d", res.Code, http.StatusBadRequest)
	}

	// test for invalid endDate
	reqBody = dtos.PostAvailabilityBody{
		StartDate: "2024-05-30",
		EndDate: "invalid",
	}
	
	jsonData, err = json.Marshal(reqBody)

	if err != nil {
		t.Log("Error:", err)
		return
	}
	req, _ = http.NewRequest("POST", "/search-availability", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	res = httptest.NewRecorder()

	reqBodyRef = &dtos.PostAvailabilityBody{}
	handlerChain = middleware.ValidateReqBody(http.HandlerFunc(Repo.SearchAvailability), reqBodyRef)

	handlerChain.ServeHTTP(res, req)

	if res.Code != http.StatusBadRequest {
		t.Errorf("SearchAvailability handler returned wrong response code for invalid endDate: got %d, wanted %d", res.Code, http.StatusBadRequest)
	}	
	
	// test for failed db search
	reqBody = dtos.PostAvailabilityBody{
		StartDate: "1955-05-30",
		EndDate: "2024-06-06",
	}
	
	jsonData, err = json.Marshal(reqBody)

	if err != nil {
		t.Log("Error:", err)
		return
	}
	req, _ = http.NewRequest("POST", "/search-availability", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	res = httptest.NewRecorder()

	reqBodyRef = &dtos.PostAvailabilityBody{}
	handlerChain = middleware.ValidateReqBody(http.HandlerFunc(Repo.SearchAvailability), reqBodyRef)

	handlerChain.ServeHTTP(res, req)

	if res.Code != http.StatusNotFound {
		t.Errorf("SearchAvailability handler returned wrong response code for failed db search: got %d, wanted %d", res.Code, http.StatusNotFound)
	}
}

func TestRepository_SearchAvailabilityByRoomId(t *testing.T){
	reqBody := dtos.PostAvailabilityBody{
		StartDate: "2024-05-30",
		EndDate: "2024-06-06",
	}
	
	jsonData, err := json.Marshal(reqBody)
    if err != nil {
        t.Log("Error:", err)
        return
    }

	// Test to make sure that a post handler is being called
	req, _ := http.NewRequest("PUT", "/search-availability/1", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	res := httptest.NewRecorder()

	handler := middleware.ValidateReqBody(http.HandlerFunc(Repo.SearchAvailabilityByRoomId), &dtos.PostAvailabilityBody{})

	handler.ServeHTTP(res, req)

	if res.Code != http.StatusMethodNotAllowed {
		t.Errorf("SearchAvailabilityByRoomId handler returned wrong response code for wrong http method: got %d, wanted %d", res.Code, http.StatusMethodNotAllowed)
	}

	// Test for invalid id in the path variable
	req, _ = http.NewRequest("POST", "/search-availability/one", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.RequestURI = "/search-availability/one"

	res = httptest.NewRecorder()

	handler = middleware.ValidateReqBody(http.HandlerFunc(Repo.SearchAvailabilityByRoomId), &dtos.PostAvailabilityBody{})

	handler.ServeHTTP(res, req)

	if res.Code != http.StatusInternalServerError {
		t.Errorf("SearchAvailabilityByRoomId handler returned wrong response code for invalid id in the url: got %d, wanted %d", res.Code, http.StatusInternalServerError)
	}
	
	// test for the right request body
	req, _ = http.NewRequest("POST", "/search-availability/1", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.RequestURI = "/search-availability/1"

	res = httptest.NewRecorder()

	handler = middleware.ValidateReqBody(http.HandlerFunc(Repo.SearchAvailabilityByRoomId), &dtos.PostAvailabilityBody{})

	handler.ServeHTTP(res, req)

	if res.Code != http.StatusFound {
		t.Errorf("SearchAvailabilityByRoomId handler returned wrong response code: got %d, wanted %d", res.Code, http.StatusFound)
	}

	// test for missing request body
	req, _ = http.NewRequest("POST", "/search-availability/1", bytes.NewBuffer([]byte(``)))
	req.Header.Set("Content-Type", "application/json")
	req.RequestURI = "/search-availability/1"

	res = httptest.NewRecorder()

	handler = middleware.ValidateReqBody(http.HandlerFunc(Repo.SearchAvailabilityByRoomId), &dtos.PostAvailabilityBody{})

	handler.ServeHTTP(res, req)

	if res.Code != http.StatusBadRequest {
		t.Errorf("SearchAvailabilityByRoomId handler returned wrong response code for missing request body: got %d, wanted %d", res.Code, http.StatusBadRequest)
	}

	// test for missing request body data in the context
	req, _ = http.NewRequest("POST", "/search-availability/1", bytes.NewBuffer([]byte(``)))
	req.Header.Set("Content-Type", "application/json")
	req.RequestURI = "/search-availability/1"

	res = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.SearchAvailabilityByRoomId)

	handler.ServeHTTP(res, req)

	if res.Code != http.StatusBadRequest {
		t.Errorf("SearchAvailabilityByRoomId handler returned wrong response code for missing request body data in the request context: got %d, wanted %d", res.Code, http.StatusBadRequest)
	}

	// test for invalid startDate
	reqBody = dtos.PostAvailabilityBody{
		StartDate: "invalid",
		EndDate: "2024-06-06",
	}
	
	jsonData, err = json.Marshal(reqBody)

	if err != nil {
		t.Log("Error:", err)
		return
	}
	req, _ = http.NewRequest("POST", "/search-availability/1", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.RequestURI = "/search-availability/1"

	res = httptest.NewRecorder()

	handler = middleware.ValidateReqBody(http.HandlerFunc(Repo.SearchAvailabilityByRoomId), &dtos.PostAvailabilityBody{})

	handler.ServeHTTP(res, req)

	if res.Code != http.StatusBadRequest {
		t.Errorf("SearchAvailabilityByRoomId handler returned wrong response code for invalid startDate: got %d, wanted %d", res.Code, http.StatusBadRequest)
	}

	// test for invalid endDate
	reqBody = dtos.PostAvailabilityBody{
		StartDate: "2024-05-30",
		EndDate: "invalid",
	}
	
	jsonData, err = json.Marshal(reqBody)

	if err != nil {
		t.Log("Error:", err)
		return
	}
	req, _ = http.NewRequest("POST", "/search-availability/1", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.RequestURI = "/search-availability/1"

	res = httptest.NewRecorder()

	handler = middleware.ValidateReqBody(http.HandlerFunc(Repo.SearchAvailabilityByRoomId), &dtos.PostAvailabilityBody{})

	handler.ServeHTTP(res, req)

	if res.Code != http.StatusBadRequest {
		t.Errorf("SearchAvailabilityByRoomId handler returned wrong response code for invalid endDate: got %d, wanted %d", res.Code, http.StatusBadRequest)
	}

	// test for failed db search
	reqBody = dtos.PostAvailabilityBody{
		StartDate: "2024-05-30",
		EndDate: "2024-06-06",
	}
	
	jsonData, err = json.Marshal(reqBody)
    if err != nil {
        t.Log("Error:", err)
        return
    }

	req, _ = http.NewRequest("POST", "/search-availability/2", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.RequestURI = "/search-availability/2"

	res = httptest.NewRecorder()

	handler = middleware.ValidateReqBody(http.HandlerFunc(Repo.SearchAvailabilityByRoomId), &dtos.PostAvailabilityBody{})

	handler.ServeHTTP(res, req)

	if res.Code != http.StatusNotFound {
		t.Errorf("SearchAvailabilityByRoomId handler returned wrong response code for failed db search: got %d, wanted %d", res.Code, http.StatusNotFound)
	}
}

func TestRepository_GetAllRooms(t *testing.T) {
	// test OK
	req, _ := http.NewRequest("GET", "/room", nil)
	req.Header.Set("Content-Type", "application/json")

	res := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.GetAllRooms)

	handler.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Errorf("GetAllRooms handler returned wrong response code: got %d, wanted %d", res.Code, http.StatusOK)
	}

	// test valid id in query param
	req, _ = http.NewRequest("GET", "/room", nil)
	req.Header.Set("Content-Type", "application/json")
	params := url.Values{}
    params.Add("id", "one")
    req.URL.RawQuery = params.Encode()

	res = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.GetAllRooms)

	handler.ServeHTTP(res, req)

	if res.Code != http.StatusNotFound {
		t.Errorf("GetAllRooms handler returned wrong response code for invalid query param 'id': got %d, wanted %d", res.Code, http.StatusNotFound)
	}

	// test for failed db operation
	req, _ = http.NewRequest("GET", "/room", nil)
	req.Header.Set("Content-Type", "application/json")	
	params = url.Values{}
    params.Add("id", "2")
    req.URL.RawQuery = params.Encode()

	res = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.GetAllRooms)

	handler.ServeHTTP(res, req)

	if res.Code != http.StatusNotFound {
		t.Errorf("GetAllRooms handler returned wrong response code for invalid query param 'id': got %d, wanted %d", res.Code, http.StatusNotFound)
	}
}

func TestRepository_GetARoomById(t *testing.T) {
	// test OK
	req, _ := http.NewRequest("GET", "/room/1", nil)
	req.Header.Set("Content-Type", "application/json")
	req.RequestURI = "/room/1"

	res := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.GetRoomById)

	handler.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Errorf("GetRoomById handler returned wrong response code: got %d, wanted %d", res.Code, http.StatusOK)
	}

	// test valid id in the path variable
	req, _ = http.NewRequest("GET", "/room", nil)
	req.Header.Set("Content-Type", "application/json")
	req.RequestURI = "/room/one"

	res = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.GetRoomById)

	handler.ServeHTTP(res, req)

	if res.Code != http.StatusInternalServerError {
		t.Errorf("GetRoomById handler returned wrong response code for invalid query param 'id': got %d, wanted %d", res.Code, http.StatusInternalServerError)
	}

	// test for failed db operation
	req, _ = http.NewRequest("GET", "/room", nil)
	req.Header.Set("Content-Type", "application/json")	
	req.RequestURI = "/room/1000"


	res = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.GetRoomById)

	handler.ServeHTTP(res, req)

	if res.Code != http.StatusNotFound {
		t.Errorf("GetRoomById handler returned wrong response code for invalid query param 'id': got %d, wanted %d", res.Code, http.StatusNotFound)
	}
}

func TestLoginHandler(t *testing.T){
	reqBody := dtos.UserLoginBody{
		Email: "johndoe@gmail.com",
	}
	
	jsonData, err := json.Marshal(reqBody)
    if err != nil {
        t.Log("Error:", err)
        return
    }

	// Test for missing request body
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer([]byte(``)))
	req.Header.Set("Content-Type", "application/json")

	res := httptest.NewRecorder()

	reqBodyRef := &dtos.UserLoginBody{}
	handler := http.HandlerFunc(Repo.LoginUser)

	handler.ServeHTTP(res, req)

	if res.Code != http.StatusBadRequest {
		t.Errorf("Login handler returned wrong response code for missing request body: got %d, wanted %d", res.Code, http.StatusBadRequest)
	}

	// Test for success
	req, _ = http.NewRequest("POST", "/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	res = httptest.NewRecorder()

	handlerChain := middleware.ValidateReqBody(http.HandlerFunc(Repo.LoginUser), reqBodyRef)

	handlerChain.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Errorf("Login handler returned wrong response code: got %d, wanted %d", res.Code, http.StatusOK)
	}
}

func TestProtectedRouteHandler(t *testing.T){
	req, _ := http.NewRequest("GET", "/protected=route", bytes.NewBuffer([]byte(``)))
	req.Header.Set("Content-Type", "application/json")

	res := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.ProtectedRoute)

	handler.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Errorf("ProtectedRoute handler returned wrong response code: got %d, wanted %d", res.Code, http.StatusOK)
	}
}