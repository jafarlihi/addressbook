package service

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/jafarlihi/addressbook/config"
	"strings"
)

func ParseAuthorizationHeader(header string) (float64, error) {
	tokenFields := strings.Fields(header)
	if len(tokenFields) != 2 {
		return 0, errors.New("Token is missing")
	}
	tokenString := tokenFields[1]
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.Config.Jwt.SigningSecret), nil
	})
	if err != nil {
		return 0, errors.New("Failed to parse the token")
	}
	var userID float64
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID = claims["userID"].(float64)
	} else {
		return 0, errors.New("Invalid token")
	}
	return userID, nil
}
