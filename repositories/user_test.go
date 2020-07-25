package repositories_test

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jafarlihi/addressbook/repositories"
	"testing"
)

func TestGetUserByUsername(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	var id uint32
	id = 1
	username := "username"
	email := "user@email.com"
	password := "password"

	rows := sqlmock.NewRows([]string{"id", "username", "email", "password"}).AddRow(id, username, email, password)
	mock.ExpectQuery("^SELECT (.+) FROM users WHERE").WithArgs(username).WillReturnRows(rows)

	user, err := repositories.GetUserByUsername(db, username)
	if err != nil {
		t.Errorf("Error was not expected while fetching the user: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}

	if user.Email != email || user.Username != username || user.ID != id || user.Password != password {
		t.Errorf("Returned user object does not match the expectations")
	}
}

func TestGetUserByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	var id uint32
	id = 1
	username := "username"
	email := "user@email.com"
	password := "password"

	rows := sqlmock.NewRows([]string{"id", "username", "email", "password"}).AddRow(id, username, email, password)
	mock.ExpectQuery("^SELECT (.+) FROM users WHERE").WithArgs(email).WillReturnRows(rows)

	user, err := repositories.GetUserByEmail(db, email)
	if err != nil {
		t.Errorf("Error was not expected while fetching the user: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}

	if user.Email != email || user.Username != username || user.ID != id || user.Password != password {
		t.Errorf("Returned user object does not match the expectations")
	}
}

func TestCreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	var id uint32
	id = 1
	username := "username"
	email := "user@email.com"
	password := "password"

	rows := sqlmock.NewRows([]string{"id"}).AddRow(id)
	mock.ExpectQuery("^INSERT INTO users").WithArgs(username, email, password).WillReturnRows(rows)

	returnedID, err := repositories.CreateUser(db, username, email, password)
	if err != nil {
		t.Errorf("Error was not expected while fetching the user: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}

	if uint32(returnedID) != id {
		t.Errorf("Returned ID '%d' does not match the expectations", returnedID)
	}
}
