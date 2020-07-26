package handlers

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"regexp"

	"github.com/dgrijalva/jwt-go"
	"github.com/jafarlihi/addressbook/config"
	"github.com/jafarlihi/addressbook/database"
	"github.com/jafarlihi/addressbook/models"
	"github.com/jafarlihi/addressbook/repositories"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(w http.ResponseWriter, r *http.Request, body Request) {
	if body.Username == "" || body.Email == "" || body.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, `{"error": "Username, email, or password field(s) is/are missing"}`)
		return
	}

	if !regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`).MatchString(body.Email) {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, `{"error": "Provided email address is malformed"}`)
		return
	}

	if len(body.Password) < 6 {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, `{"error": "Password length can't be smaller than 6"}`)
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, `{"error": "Failed to hash the password"}`)
		return
	}

	id, err := repositories.CreateUser(database.Database, body.Username, body.Email, string(passwordHash))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, `{"error": "Failed to create the user, it might already exist"}`)
		return
	}

	jsonResponse, _ := json.Marshal(map[string]int64{"id": id})
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, string(jsonResponse))
}

func CreateToken(w http.ResponseWriter, r *http.Request, body Request) {
	if body.Username == "" && body.Email == "" {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, `{"error": "Username and email fields are missing, at least one is required"}`)
		return
	}

	if body.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, `{"error": "Password field is missing"}`)
		return
	}

	var err error
	var user *models.User
	if body.Username != "" {
		user, err = repositories.GetUserByUsername(database.Database, body.Username)
	} else {
		user, err = repositories.GetUserByEmail(database.Database, body.Email)
	}
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, `{"error": "User does not exist"}`)
		default:
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, `{"error": "Failed to get the user"}`)
		}
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, `{"error": "Wrong password"}`)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": user.ID,
	})
	tokenString, err := token.SignedString([]byte(config.Config.Jwt.SigningSecret))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, `{"error": "Failed to create token"}`)
		return
	}

	response := struct {
		Token string      `json:"token"`
		User  models.User `json:"user"`
	}{
		tokenString,
		*user,
	}
	response.User.Password = ""
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, `{"error": "Failed to marshal the response to JSON"}`)
		return
	}
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, string(jsonResponse))
}
