package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jafarlihi/addressbook/database"
	"github.com/jafarlihi/addressbook/repositories"
)

type contactCreationRequest struct {
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Email   string `json:"email"`
}

func CreateContact(w http.ResponseWriter, r *http.Request, userID uint32) {
	var ccr contactCreationRequest
	err := json.NewDecoder(r.Body).Decode(&ccr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, `{"error": "Request body couldn't be parsed as JSON"}`)
		return
	}

	if ccr.Name == "" || ccr.Surname == "" || ccr.Email == "" {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, `{"error": "Name, surname, and/or email field(s) is/are missing"}`)
		return
	}

	if !regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`).MatchString(ccr.Email) {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, `{"error": "Provided email address is malformed"}`)
		return
	}

	id, err := repositories.CreateContact(database.Database, userID, ccr.Name, ccr.Surname, ccr.Email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, `{"error": "Failed to create the contact"}`)
		return
	}

	jsonResponse, _ := json.Marshal(map[string]int64{"id": id})
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, string(jsonResponse))
}

func DeleteContact(w http.ResponseWriter, r *http.Request, userID uint32) {
	params := mux.Vars(r)
	idString := params["id"]
	id, err := strconv.ParseUint(idString, 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, `{"error": "Provided ID can't be parsed as an integer"}`)
		return
	}

	contact, err := repositories.GetContact(database.Database, uint32(id))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, `{"error": "Requested contact does not exist"}`)
		return
	}

	if contact.UserID != userID {
		w.WriteHeader(http.StatusUnauthorized)
		io.WriteString(w, `{"error": "Can't delete contact belonging to another user"}`)
		return
	}

	err = repositories.DeleteContact(database.Database, uint32(id))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, `{"error": "Failed to delete the contact"}`)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func GetContacts(w http.ResponseWriter, r *http.Request, userID uint32) {
	contacts, err := repositories.GetContactsByUserID(database.Database, userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, `{"error": "Failed to get the contacts"}`)
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

func GetContact(w http.ResponseWriter, r *http.Request, userID uint32) {
	params := mux.Vars(r)
	idString := params["id"]
	id, err := strconv.ParseUint(idString, 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, `{"error": "Provided ID can't be parsed as an integer"}`)
		return
	}

	contact, err := repositories.GetContact(database.Database, uint32(id))
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

	jsonResponse, err := json.Marshal(contact)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, `{"error": "Failed to marshal the result to JSON"}`)
		return
	}

	w.WriteHeader(http.StatusOK)
	io.WriteString(w, string(jsonResponse))
}
