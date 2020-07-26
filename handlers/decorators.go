package handlers

import (
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
