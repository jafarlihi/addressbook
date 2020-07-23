package main

import (
	gorillaHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jafarlihi/addressbook/config"
	"github.com/jafarlihi/addressbook/database"
	"github.com/jafarlihi/addressbook/handlers"
	"github.com/jafarlihi/addressbook/logger"
	"net/http"
)

func main() {
	logger.InitLogger()
	config.InitConfig()
	database.InitDatabase()

	router := mux.NewRouter()
	router.HandleFunc("/api/user", handlers.CreateUser).Methods("POST")
	router.HandleFunc("/api/user/token", handlers.CreateToken).Methods("POST")
	router.HandleFunc("/api/contact", handlers.CreateContact).Methods("POST")
	router.HandleFunc("/api/contact/{id}", handlers.DeleteContact).Methods("DELETE")
	router.HandleFunc("/api/contact", handlers.GetContacts).Methods("GET")
	router.HandleFunc("/api/contact/{id}", handlers.GetContact).Methods("GET")
	router.HandleFunc("/api/contact-list", handlers.CreateContactList).Methods("POST")
	router.HandleFunc("/api/contact-list/{id}", handlers.DeleteContactList).Methods("DELETE")
	router.HandleFunc("/api/contact-list", handlers.GetContactLists).Methods("GET")
	router.HandleFunc("/api/contact-list/{id}", handlers.GetContactList).Methods("GET")
	router.HandleFunc("/api/contact-list/search", handlers.SearchContactLists).Methods("POST")
	router.HandleFunc("/api/contact-list/{id}/contact", handlers.ListContactsOfContactList).Methods("GET")
	router.HandleFunc("/api/contact-list/{id}/contact", handlers.AddToContactList).Methods("POST")
	router.HandleFunc("/api/contact-list/{id}/contact", handlers.RemoveFromContactList).Methods("DELETE")

	origins := gorillaHandlers.AllowedOrigins([]string{"*"})
	headers := gorillaHandlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := gorillaHandlers.AllowedMethods([]string{"GET", "POST", "DELETE"})

	logger.Log.Info("Starting HTTP server listening at " + config.Config.HttpServer.Port)
	logger.Log.Critical(http.ListenAndServe(":"+config.Config.HttpServer.Port, gorillaHandlers.CORS(origins, headers, methods)(router)))
}
