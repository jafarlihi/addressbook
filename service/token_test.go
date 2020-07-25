package service_test

import (
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/jafarlihi/addressbook/config"
	"github.com/jafarlihi/addressbook/service"
)

func TestParseAuthorizationHeader(t *testing.T) {
	jwtSecret := "secret"
	config.Config.Jwt.SigningSecret = jwtSecret

	var id uint32
	id = 1

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": id,
	})
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		t.Fatal("Failed to create JWT token")
	}

	header := "Bearer " + tokenString

	returnedID, err := service.ParseAuthorizationHeader(header)
	if err != nil {
		t.Errorf("ParseAuthorizationHeader returned error %s", err.Error())
	}

	if returnedID != id {
		t.Errorf("Returned ID %d does not match expected ID %d", returnedID, id)
	}
}
