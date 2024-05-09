package main

import (
	"Aitunder/mongodb"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAddUser(t *testing.T) {
	requestBody := `{"name": "Kabob","email": "test@example.com", "password": "Test123!"}`
	req, err := http.NewRequest("POST", "/api/signUp", strings.NewReader(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	mongodb.AddUser(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code. Expected %v, got %v", http.StatusOK, status)
	}

	expected := `{"message":"User registered successfully. Please verify your email.","status":200}`
	actual := strings.TrimSpace(rr.Body.String())

	if actual != expected {
		t.Errorf("Handler returned unexpected body. Expected %v, got %v", expected, actual)
	}
}
