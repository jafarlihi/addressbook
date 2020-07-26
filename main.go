package main

import (
	"net/http"

	gorillaHandlers "github.com/gorilla/handlers"
	"github.com/jafarlihi/addressbook/config"
	"github.com/jafarlihi/addressbook/database"
	"github.com/jafarlihi/addressbook/logger"
	"github.com/jafarlihi/addressbook/router"
)

func main() {
	logger.InitLogger()
	config.InitConfig()
	database.InitDatabase()

	router := router.ConstructRouter()

	origins := gorillaHandlers.AllowedOrigins([]string{"*"})
	headers := gorillaHandlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := gorillaHandlers.AllowedMethods([]string{"GET", "POST", "DELETE"})

	logger.Log.Info("Starting HTTP server listening at " + config.Config.HttpServer.Port)
	logger.Log.Critical(http.ListenAndServe(":"+config.Config.HttpServer.Port, gorillaHandlers.CORS(origins, headers, methods)(router)))
}
