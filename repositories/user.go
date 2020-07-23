package repositories

import (
	"github.com/jafarlihi/addressbook/database"
	"github.com/jafarlihi/addressbook/logger"
	"github.com/jafarlihi/addressbook/models"
)

func GetUserByUsername(username string) (*models.User, error) {
	sql := "SELECT id, username, email, password FROM users WHERE username = $1"
	row := database.Database.QueryRow(sql, username)
	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		logger.Log.Error("Failed to SELECT a user, error: " + err.Error())
		return nil, err
	}
	return &user, nil
}

func GetUserByEmail(email string) (*models.User, error) {
	sql := "SELECT id, username, email, password FROM users WHERE email = $1"
	row := database.Database.QueryRow(sql, email)
	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		logger.Log.Error("Failed to SELECT a user, error: " + err.Error())
		return nil, err
	}
	return &user, nil
}

func CreateUser(username string, email string, password string) (int64, error) {
	sql := "INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id"
	var id int64
	err := database.Database.QueryRow(sql, username, email, password).Scan(&id)
	if err != nil {
		logger.Log.Error("Failed to INSERT a new user, error: " + err.Error())
		return 0, err
	}
	return id, nil
}
