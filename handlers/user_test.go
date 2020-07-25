package handlers_test

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/dgrijalva/jwt-go"
	"github.com/jafarlihi/addressbook/config"
	"github.com/jafarlihi/addressbook/database"
	"github.com/jafarlihi/addressbook/handlers"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateUserWithNoBody(t *testing.T) {
	req, err := http.NewRequest("POST", "/api/user", strings.NewReader(""))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.CreateUser)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	expected := `{"error": "Request body couldn't be parsed as JSON"}`
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestCreateUserWithMissingField(t *testing.T) {
	req, err := http.NewRequest("POST", "/api/user", strings.NewReader(`{"username": "user"}`))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.CreateUser)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	expected := `{"error": "Username, email, or password field(s) is/are missing"}`
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestCreateUserWithInvalidEmail(t *testing.T) {
	req, err := http.NewRequest("POST", "/api/user", strings.NewReader(`{"username": "user", "password": "pass", "email": "invalid"}`))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.CreateUser)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	expected := `{"error": "Provided email address is malformed"}`
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestCreateUserWithShortPassword(t *testing.T) {
	req, err := http.NewRequest("POST", "/api/user", strings.NewReader(`{"username": "user", "password": "pass", "email": "valid@mail.com"}`))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.CreateUser)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	expected := `{"error": "Password length can't be smaller than 6"}`
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestCreateUser(t *testing.T) {
	username := "user"
	email := "valid@mail.com"
	password := "password"

	req, err := http.NewRequest("POST", "/api/user", strings.NewReader(`{"username": "`+username+`", "password": "`+password+`", "email": "`+email+`"}`))
	if err != nil {
		t.Fatal(err)
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	database.Database = db

	var id uint32
	id = 1

	rows := sqlmock.NewRows([]string{"id"}).AddRow(id)
	mock.ExpectQuery("^INSERT INTO users").WithArgs(username, email, sqlmock.AnyArg()).WillReturnRows(rows)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.CreateUser)
	handler.ServeHTTP(rr, req)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"id":` + fmt.Sprint(id) + `}`
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestCreateToken(t *testing.T) {
	username := "user"
	password := "password"

	req, err := http.NewRequest("POST", "/api/user/token", strings.NewReader(`{"username": "`+username+`", "password": "`+password+`"}`))
	if err != nil {
		t.Fatal(err)
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	database.Database = db

	var id uint32
	id = 1
	email := "valid@mail.com"
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatal("Failed to hash password")
	}

	rows := sqlmock.NewRows([]string{"id", "username", "email", "password"}).AddRow(id, username, email, string(passwordHash))
	mock.ExpectQuery("^SELECT (.*) FROM users").WithArgs(username).WillReturnRows(rows)

	jwtSecret := "secret"
	config.Config.Jwt.SigningSecret = jwtSecret

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": id,
	})
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		t.Fatal("Failed to create JWT token")
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.CreateToken)
	handler.ServeHTTP(rr, req)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"token":"` + tokenString + `","user":{"id":`+fmt.Sprint(id)+`,"username":"`+username+`","email":"`+email+`","password":""}}`
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
