package repositories_test

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jafarlihi/addressbook/repositories"
	"testing"
)

func TestCreateContactList(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	var id uint32
	id = 1
	var userID uint32
	userID = 1
	name := "name"

	rows := sqlmock.NewRows([]string{"id"}).AddRow(id)
	mock.ExpectQuery("^INSERT INTO contact_lists").WithArgs(userID, name).WillReturnRows(rows)

	returnedID, err := repositories.CreateContactList(db, userID, name)
	if err != nil {
		t.Errorf("Error was not expected while creating the contact-list: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}

	if uint32(returnedID) != id {
		t.Errorf("Returned ID '%d' does not match the expectations", returnedID)
	}
}
