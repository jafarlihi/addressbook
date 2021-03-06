package database

import (
	"database/sql"
	"io/ioutil"
	"os"

	"github.com/jafarlihi/addressbook/config"
	"github.com/jafarlihi/addressbook/logger"
	_ "github.com/lib/pq"
)

var Database *sql.DB

func InitDatabase() {
	var err error
	Database, err = sql.Open("postgres", config.Config.Database.Url)
	if err != nil {
		logger.Log.Error("Failed to connect to the database, error: " + err.Error())
		os.Exit(1)
	}

	err = Database.Ping()
	if err != nil {
		logger.Log.Error("Failed to connect to the database, error: " + err.Error())
		os.Exit(1)
	}

	schemaBytes, err := ioutil.ReadFile("schema.sql")
	if err != nil {
		logger.Log.Warningf("Failed to read the schema.sql for database schema initialization. Skipping procedure. Error: " + err.Error())
		return
	}
	_, err = Database.Exec(string(schemaBytes))
	if err != nil {
		logger.Log.Warningf("Failed to (re-)initialize the schema, error: " + err.Error())
	}
}
