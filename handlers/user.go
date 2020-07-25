package handlers

import (
	"database/sql"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/jafarlihi/addressbook/config"
	"github.com/jafarlihi/addressbook/database"
	"github.com/jafarlihi/addressbook/models"
	"github.com/jafarlihi/addressbook/repositories"
	"golang.org/x/crypto/bcrypt"
	"io"
	"net/http"
	"regexp"
)

type accountCreationRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var acr accountCreationRequest
	err := json.NewDecoder(r.Body).Decode(&acr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, `{"error": "Request body couldn't be parsed as JSON"}`)
		return
	}
	if acr.Username == "" || acr.Email == "" || acr.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, `{"error": "Username, email, or password field(s) is/are missing"}`)
		return
	}
	if !regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`).MatchString(acr.Email) {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, `{"error": "Provided email address is malformed"}`)
		return
	}
	if len(acr.Password) < 6 {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, `{"error": "Password length can't be smaller than 6"}`)
		return
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(acr.Password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, `{"error": "Failed to hash the password"}`)
		return
	}
	id, err := repositories.CreateUser(database.Database, acr.Username, acr.Email, string(passwordHash))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, `{"error": "Failed to create the user"}`)
		return
	}
	jsonResponse, _ := json.Marshal(map[string]int64{"id": id})
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, string(jsonResponse))
}

type tokenCreationRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type tokenCreationResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

func CreateToken(w http.ResponseWriter, r *http.Request) {
	var tcr tokenCreationRequest
	err := json.NewDecoder(r.Body).Decode(&tcr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, `{"error": "Request body couldn't be parsed as JSON"}`)
		return
	}
	if tcr.Username == "" && tcr.Email == "" {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, `{"error": "Username and email fields are missing, at least one is required"}`)
		return
	}
	if tcr.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, `{"error": "Password field is missing"}`)
		return
	}
	var user *models.User
	if tcr.Username != "" {
		user, err = repositories.GetUserByUsername(database.Database, tcr.Username)
	} else {
		user, err = repositories.GetUserByEmail(database.Database, tcr.Email)
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
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(tcr.Password))
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
	var response tokenCreationResponse
	response.Token = tokenString
	response.User = *user
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
