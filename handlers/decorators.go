package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/jafarlihi/addressbook/services"
)

func Authenticated(w http.ResponseWriter, r *http.Request, f func(http.ResponseWriter, *http.Request, uint32)) {
	userID, err := services.ParseAuthorizationHeader(r.Header.Get("Authorization"))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		io.WriteString(w, `{"error": "`+err.Error()+`"}`)
		return
	}

	f(w, r, userID)
}

func AuthenticatedWithRequestBody(w http.ResponseWriter, r *http.Request, f func(http.ResponseWriter, *http.Request, uint32, Request)) {
	var body Request
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, `{"error": "Request body couldn't be parsed as JSON"}`)
		return
	}

	userID, err := services.ParseAuthorizationHeader(r.Header.Get("Authorization"))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		io.WriteString(w, `{"error": "`+err.Error()+`"}`)
		return
	}

	f(w, r, userID, body)
}
