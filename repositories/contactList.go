package repositories

import (
	"github.com/jafarlihi/addressbook/database"
	"github.com/jafarlihi/addressbook/logger"
	"github.com/jafarlihi/addressbook/models"
)

func CreateContactList(userID uint32, name string) (int64, error) {
	sql := "INSERT INTO contact_lists (user_id, name) VALUES ($1, $2) RETURNING id"
	var id int64
	err := database.Database.QueryRow(sql, userID, name).Scan(&id)
	if err != nil {
		logger.Log.Error("Failed to INSERT a new contact-list, error: " + err.Error())
		return 0, err
	}
	return id, nil
}

func GetContactList(id uint32) (*models.ContactList, error) {
	sql := "SELECT id, user_id, name FROM contact_lists WHERE id = $1"
	row := database.Database.QueryRow(sql, id)
	var contactList models.ContactList
	err := row.Scan(&contactList.ID, &contactList.UserID, &contactList.Name)
	if err != nil {
		logger.Log.Error("Failed to SELECT a contact, error: " + err.Error())
		return nil, err
	}
	return &contactList, nil
}

func DeleteContactList(id uint32) error {
	sql := "DELETE FROM contact_lists WHERE id = $1"
	_, err := database.Database.Exec(sql, id)
	if err != nil {
		logger.Log.Error("Failed to DELETE a contact-list, error: " + err.Error())
		return err
	}
	return nil
}

func GetContactListsByUserID(userID uint32) ([]*models.ContactList, error) {
	sql := "SELECT id, user_id, name FROM contact_lists WHERE user_id = $1"
	rows, err := database.Database.Query(sql, userID)
	if err != nil {
		logger.Log.Error("Failed to SELECT contact-lists, error: " + err.Error())
		return nil, err
	}
	defer rows.Close()
	contactLists := make([]*models.ContactList, 0)
	for rows.Next() {
		contactList := &models.ContactList{}
		if err := rows.Scan(&contactList.ID, &contactList.UserID, &contactList.Name); err != nil {
			logger.Log.Error("Failed to scan SELECTed row of contact-lists, error: " + err.Error())
			return nil, err
		}
		contactLists = append(contactLists, contactList)
	}
	return contactLists, nil
}

func SearchContactListsByName(userID uint32, term string) ([]*models.ContactList, error) {
	sql := "SELECT id, user_id, name FROM contact_lists WHERE user_id = $1 AND name ILIKE '%' || $2 || '%'"
	rows, err := database.Database.Query(sql, userID, term)
	if err != nil {
		logger.Log.Error("Failed to SELECT contact-lists, error: " + err.Error())
		return nil, err
	}
	defer rows.Close()
	contactLists := make([]*models.ContactList, 0)
	for rows.Next() {
		contactList := &models.ContactList{}
		if err := rows.Scan(&contactList.ID, &contactList.UserID, &contactList.Name); err != nil {
			logger.Log.Error("Failed to scan SELECTed row of contact-lists, error: " + err.Error())
			return nil, err
		}
		contactLists = append(contactLists, contactList)
	}
	return contactLists, nil
}

func AddContactToContactList(contactListID uint32, contactID uint32) error {
	sql := "INSERT INTO contact_list_entries (contact_list, contact) VALUES ($1, $2)"
	_, err := database.Database.Exec(sql, contactListID, contactID)
	if err != nil {
		logger.Log.Error("Failed to INSERT a new contact-list-entry, error: " + err.Error())
		return err
	}
	return nil
}

func DeleteContactFromContactList(contactListID uint32, contactID uint32) error {
	sql := "DELETE FROM contact_list_entries WHERE contact_list = $1 AND contact = $2"
	_, err := database.Database.Exec(sql, contactListID, contactID)
	if err != nil {
		logger.Log.Error("Failed to DELETE a contact-list-entry, error: " + err.Error())
		return err
	}
	return nil
}
