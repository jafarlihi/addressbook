package handlers_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/jafarlihi/addressbook/config"
	"github.com/jafarlihi/addressbook/database"
	"github.com/jafarlihi/addressbook/handlers"
)

func TestCreateContactListNoBody(t *testing.T) {
	req, err := http.NewRequest("POST", "/api/contact-list", strings.NewReader(""))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.CreateContactList)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	expected := `{"error": "Request body couldn't be parsed as JSON"}`
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestCreateContactListWithNoNameField(t *testing.T) {
	req, err := http.NewRequest("POST", "/api/contact-list", strings.NewReader(`{"notName": "something"}`))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.CreateContactList)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	expected := `{"error": "Name field is missing"}`
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestCreateContactListWithNoToken(t *testing.T) {
	req, err := http.NewRequest("POST", "/api/contact-list", strings.NewReader(`{"name": "something"}`))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.CreateContactList)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
	}

	expected := `{"error": "Token is missing"}`
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestCreateContactList(t *testing.T) {
	jwtSecret := "secret"
	config.Config.Jwt.SigningSecret = jwtSecret

	var userID uint32
	userID = 1
	var contactListID uint32
	contactListID = 2
	name := "name"

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

	rows := sqlmock.NewRows([]string{"id"}).AddRow(contactListID)
	mock.ExpectQuery("^INSERT INTO contact_lists").WithArgs(userID, name).WillReturnRows(rows)

	req, err := http.NewRequest("POST", "/api/contact-list", strings.NewReader(`{"name": "`+name+`"}`))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("Authorization", "Bearer "+tokenString)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.CreateContactList)
	handler.ServeHTTP(rr, req)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"id":` + fmt.Sprint(contactListID) + `}`
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestDeleteContactListWithNoToken(t *testing.T) {
	req, err := http.NewRequest("DELETE", "/api/contact-list/1", strings.NewReader(""))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/api/contact-list/{id}", handlers.DeleteContactList).Methods("DELETE")
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
	}

	expected := `{"error": "Token is missing"}`
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestDeleteContactList(t *testing.T) {
	jwtSecret := "secret"
	config.Config.Jwt.SigningSecret = jwtSecret

	var userID uint32
	userID = 1
	var contactListID uint32
	contactListID = 2
	name := "name"

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

	rows := sqlmock.NewRows([]string{"id", "user_id", "name"}).AddRow(contactListID, userID, name)
	mock.ExpectQuery("^SELECT (.*) FROM contact_lists").WithArgs(contactListID).WillReturnRows(rows)
	mock.ExpectQuery("^DELETE FROM contact_lists").WithArgs(contactListID).WillReturnRows(rows)

	req, err := http.NewRequest("DELETE", "/api/contact-list/"+fmt.Sprint(contactListID), strings.NewReader(""))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("Authorization", "Bearer "+tokenString)

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/api/contact-list/{id}", handlers.DeleteContactList).Methods("DELETE")
	router.ServeHTTP(rr, req)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := ""
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
