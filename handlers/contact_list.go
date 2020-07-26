package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jafarlihi/addressbook/database"
	"github.com/jafarlihi/addressbook/repositories"
	"github.com/jafarlihi/addressbook/services"
)

func CreateContactList(w http.ResponseWriter, r *http.Request, userID uint32, body Request) {
	if body.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, `{"error": "Name field is missing"}`)
		return
	}

	id, err := repositories.CreateContactList(database.Database, userID, body.Name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, `{"error": "Failed to create the contact-list"}`)
		return
	}

	jsonResponse, _ := json.Marshal(map[string]int64{"id": id})
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, string(jsonResponse))
}

func DeleteContactList(w http.ResponseWriter, r *http.Request, userID uint32) {
	params := mux.Vars(r)
	idString := params["id"]
	id, err := strconv.ParseUint(idString, 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, `{"error": "Provided ID can't be parsed as an integer"}`)
		return
	}

	contactList, err := repositories.GetContactList(database.Database, uint32(id))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, `{"error": "Requested contact-list does not exist"}`)
		return
	}

	if contactList.UserID != userID {
		w.WriteHeader(http.StatusUnauthorized)
		io.WriteString(w, `{"error": "Can't delete contact-list belonging to another user"}`)
		return
	}

	err = repositories.DeleteContactList(database.Database, uint32(id))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, `{"error": "Failed to delete the contact-list"}`)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func GetContactLists(w http.ResponseWriter, r *http.Request, userID uint32) {
	userID, err := services.ParseAuthorizationHeader(r.Header.Get("Authorization"))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		io.WriteString(w, `{"error": "`+err.Error()+`"}`)
		return
	}

	contactLists, err := repositories.GetContactListsByUserID(database.Database, userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, `{"error": "Failed to get the contact-lists"}`)
		return
	}

	jsonResponse, err := json.Marshal(contactLists)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, `{"error": "Failed to marshal the result to JSON"}`)
		return
	}
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, string(jsonResponse))
}

func GetContactList(w http.ResponseWriter, r *http.Request, userID uint32) {
	params := mux.Vars(r)
	idString := params["id"]
	id, err := strconv.ParseUint(idString, 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, `{"error": "Provided ID can't be parsed as an integer"}`)
		return
	}

	contactList, err := repositories.GetContactList(database.Database, uint32(id))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, `{"error": "Requested contact-list does not exist"}`)
		return
	}

	if contactList.UserID != userID {
		w.WriteHeader(http.StatusUnauthorized)
		io.WriteString(w, `{"error": "Can't fetch contact-list belonging to another user"}`)
		return
	}

	jsonResponse, err := json.Marshal(contactList)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, `{"error": "Failed to marshal the result to JSON"}`)
		return
	}
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, string(jsonResponse))
}

func SearchContactLists(w http.ResponseWriter, r *http.Request, userID uint32, body Request) {
	if body.Term == "" {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, `{"error": "Term field is missing"}`)
		return
	}

	contactLists, err := repositories.SearchContactListsByName(database.Database, userID, body.Term)

	jsonResponse, err := json.Marshal(contactLists)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, `{"error": "Failed to marshal the result to JSON"}`)
		return
	}
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, string(jsonResponse))
}

func GetContactsOfContactList(w http.ResponseWriter, r *http.Request, userID uint32) {
	params := mux.Vars(r)
	idString := params["id"]
	id, err := strconv.ParseUint(idString, 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, `{"error": "Provided ID can't be parsed as an integer"}`)
		return
	}

	contactList, err := repositories.GetContactList(database.Database, uint32(id))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, `{"error": "Requested contact-list does not exist"}`)
		return
	}

	if contactList.UserID != userID {
		w.WriteHeader(http.StatusUnauthorized)
		io.WriteString(w, `{"error": "Can't fetch contact-list belonging to another user"}`)
		return
	}

	contacts, err := repositories.GetContactsOfContactList(database.Database, contactList.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, `{"error": "Failed to fetch contacts"}`)
		return
	}

	jsonResponse, err := json.Marshal(contacts)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, `{"error": "Failed to marshal the result to JSON"}`)
		return
	}
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, string(jsonResponse))
}

func AddToContactList(w http.ResponseWriter, r *http.Request, userID uint32, body Request) {
	if body.ID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, `{"error": "ID field is missing"}`)
		return
	}

	params := mux.Vars(r)
	idString := params["id"]
	id, err := strconv.ParseUint(idString, 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, `{"error": "Provided ID can't be parsed as an integer"}`)
		return
	}

	contactList, err := repositories.GetContactList(database.Database, uint32(id))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, `{"error": "Requested contact-list does not exist"}`)
		return
	}

	if contactList.UserID != userID {
		w.WriteHeader(http.StatusUnauthorized)
		io.WriteString(w, `{"error": "Can't fetch contact-list belonging to another user"}`)
		return
	}

	contact, err := repositories.GetContact(database.Database, body.ID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, `{"error": "Requested contact does not exist"}`)
		return
	}
	if contact.UserID != userID {
		w.WriteHeader(http.StatusUnauthorized)
		io.WriteString(w, `{"error": "Can't fetch contact belonging to another user"}`)
		return
	}

	err = repositories.AddContactToContactList(database.Database, contactList.ID, contact.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, `{"error": "Failed to add contact to contact-list"}`)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func RemoveFromContactList(w http.ResponseWriter, r *http.Request, userID uint32, body Request) {
	if body.ID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, `{"error": "ID field is missing"}`)
		return
	}

	params := mux.Vars(r)
	idString := params["id"]
	id, err := strconv.ParseUint(idString, 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, `{"error": "Provided ID can't be parsed as an integer"}`)
		return
	}

	contactList, err := repositories.GetContactList(database.Database, uint32(id))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, `{"error": "Requested contact-list does not exist"}`)
		return
	}

	if contactList.UserID != userID {
		w.WriteHeader(http.StatusUnauthorized)
		io.WriteString(w, `{"error": "Can't fetch contact-list belonging to another user"}`)
		return
	}

	err = repositories.DeleteContactFromContactList(database.Database, contactList.ID, body.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, `{"error": "Failed to add contact to contact-list"}`)
		return
	}

	w.WriteHeader(http.StatusOK)
}
