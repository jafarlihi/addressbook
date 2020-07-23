package repositories

import (
	"github.com/jafarlihi/addressbook/database"
	"github.com/jafarlihi/addressbook/logger"
	"github.com/jafarlihi/addressbook/models"
)

func CreateContact(userID uint32, name string, surname string, email string) (int64, error) {
	sql := "INSERT INTO contacts (user_id, name, surname, email) VALUES ($1, $2, $3, $4) RETURNING id"
	var id int64
	err := database.Database.QueryRow(sql, userID, name, surname, email).Scan(&id)
	if err != nil {
		logger.Log.Error("Failed to INSERT a new contact, error: " + err.Error())
		return 0, err
	}
	return id, nil
}

func GetContact(id uint32) (*models.Contact, error) {
	sql := "SELECT id, user_id, name, surname, email FROM contacts WHERE id = $1"
	row := database.Database.QueryRow(sql, id)
	var contact models.Contact
	err := row.Scan(&contact.ID, &contact.UserID, &contact.Name, &contact.Surname, &contact.Email)
	if err != nil {
		logger.Log.Error("Failed to SELECT a contact, error: " + err.Error())
		return nil, err
	}
	return &contact, nil
}

func DeleteContact(id uint32) error {
	sql := "DELETE FROM contacts WHERE id = $1"
	_, err := database.Database.Exec(sql, id)
	if err != nil {
		logger.Log.Error("Failed to DELETE a contact, error: " + err.Error())
		return err
	}
	return nil
}

func GetContactsByUserID(userID uint32) ([]*models.Contact, error) {
	sql := "SELECT id, user_id, name, surname, email FROM contacts WHERE user_id = $1"
	rows, err := database.Database.Query(sql, userID)
	if err != nil {
		logger.Log.Error("Failed to SELECT contacts, error: " + err.Error())
		return nil, err
	}
	defer rows.Close()
	contacts := make([]*models.Contact, 0)
	for rows.Next() {
		contact := &models.Contact{}
		if err := rows.Scan(&contact.ID, &contact.UserID, &contact.Name, &contact.Surname, &contact.Email); err != nil {
			logger.Log.Error("Failed to scan SELECTed row of contacts, error: " + err.Error())
			return nil, err
		}
		contacts = append(contacts, contact)
	}
	return contacts, nil
}

func GetContactsOfContactList(contactListID uint32) ([]*models.Contact, error) {
	sql := "SELECT id, user_id, name, surname, email FROM contacts WHERE id IN (SELECT contact FROM contact_list_entries WHERE contact_list = $1)"
	rows, err := database.Database.Query(sql, contactListID)
	if err != nil {
		logger.Log.Error("Failed to SELECT contacts, error: " + err.Error())
		return nil, err
	}
	defer rows.Close()
	contacts := make([]*models.Contact, 0)
	for rows.Next() {
		contact := &models.Contact{}
		if err := rows.Scan(&contact.ID, &contact.UserID, &contact.Name, &contact.Surname, &contact.Email); err != nil {
			logger.Log.Error("Failed to scan SELECTed row of contacts, error: " + err.Error())
			return nil, err
		}
		contacts = append(contacts, contact)
	}
	return contacts, nil
}
