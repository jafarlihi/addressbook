package handlers_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/dgrijalva/jwt-go"
	"github.com/jafarlihi/addressbook/config"
	"github.com/jafarlihi/addressbook/database"
	"github.com/jafarlihi/addressbook/handlers"
)

func TestCreateContactNoBody(t *testing.T) {
	req, err := http.NewRequest("POST", "/api/contact", strings.NewReader(""))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.CreateContact)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	expected := `{"error": "Request body couldn't be parsed as JSON"}`
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestCreateContactWithMissingField(t *testing.T) {
	req, err := http.NewRequest("POST", "/api/contact", strings.NewReader(`{"name": "name", "surname": "surname"}`))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.CreateContact)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	expected := `{"error": "Name, surname, and/or email field(s) is/are missing"}`
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestCreateContactWithNoAuthorizationHeader(t *testing.T) {
	req, err := http.NewRequest("POST", "/api/contact", strings.NewReader(`{"name": "name", "surname": "surname", "email": "email@mail.com"}`))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.CreateContact)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
	}

	expected := `{"error": "Token is missing"}`
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestCreateContact(t *testing.T) {
	jwtSecret := "secret"
	config.Config.Jwt.SigningSecret = jwtSecret

	var contactID uint32
	contactID = 2
	var userID uint32
	userID = 1
	name := "name"
	surname := "surname"
	email := "valid@mail.com"

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
	})
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		t.Fatal("Failed to create JWT token")
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	database.Database = db

	rows := sqlmock.NewRows([]string{"id"}).AddRow(contactID)
	mock.ExpectQuery("^INSERT INTO contacts").WithArgs(userID, name, surname, email).WillReturnRows(rows)

	req, err := http.NewRequest("POST", "/api/contact", strings.NewReader(`{"name": "`+name+`", "surname": "`+surname+`", "email": "`+email+`"}`))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("Authorization", "Bearer "+tokenString)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.CreateContact)
	handler.ServeHTTP(rr, req)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"id":` + fmt.Sprint(contactID) + `}`
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
