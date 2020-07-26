package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jafarlihi/addressbook/handlers"
)

func ConstructRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/user", handlers.CreateUser).Methods("POST")
	router.HandleFunc("/api/user/token", handlers.CreateToken).Methods("POST")
	router.HandleFunc("/api/contact", func(w http.ResponseWriter, r *http.Request) {
		handlers.Authenticated(w, r, handlers.CreateContact)
	}).Methods("POST")
	router.HandleFunc("/api/contact/{id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.Authenticated(w, r, handlers.DeleteContact)
	}).Methods("DELETE")
	router.HandleFunc("/api/contact", func(w http.ResponseWriter, r *http.Request) {
		handlers.Authenticated(w, r, handlers.GetContacts)
	}).Methods("GET")
	router.HandleFunc("/api/contact/{id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.Authenticated(w, r, handlers.GetContact)
	}).Methods("GET")
	router.HandleFunc("/api/contact-list", func(w http.ResponseWriter, r *http.Request) {
		handlers.Authenticated(w, r, handlers.CreateContactList)
	}).Methods("POST")
	router.HandleFunc("/api/contact-list/{id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.Authenticated(w, r, handlers.DeleteContactList)
	}).Methods("DELETE")
	router.HandleFunc("/api/contact-list", func(w http.ResponseWriter, r *http.Request) {
		handlers.Authenticated(w, r, handlers.GetContactLists)
	}).Methods("GET")
	router.HandleFunc("/api/contact-list/{id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.Authenticated(w, r, handlers.GetContactList)
	}).Methods("GET")
	router.HandleFunc("/api/contact-list/search", func(w http.ResponseWriter, r *http.Request) {
		handlers.Authenticated(w, r, handlers.SearchContactLists)
	}).Methods("POST")
	router.HandleFunc("/api/contact-list/{id}/contact", func(w http.ResponseWriter, r *http.Request) {
		handlers.Authenticated(w, r, handlers.GetContactsOfContactList)
	}).Methods("GET")
	router.HandleFunc("/api/contact-list/{id}/contact", func(w http.ResponseWriter, r *http.Request) {
		handlers.Authenticated(w, r, handlers.AddToContactList)
	}).Methods("POST")
	router.HandleFunc("/api/contact-list/{id}/contact", func(w http.ResponseWriter, r *http.Request) {
		handlers.Authenticated(w, r, handlers.RemoveFromContactList)
	}).Methods("DELETE")
	return router
}
