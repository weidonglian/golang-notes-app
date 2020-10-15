package util

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/jwtauth"
	"net/http"
)

func getClaimsFromRequest(r *http.Request) jwt.MapClaims {
	if _, claims, err := jwtauth.FromContext(r.Context()); err != nil {
		panic(err)
	} else {
		return claims
	}
}

func GetUserIDFromRequest(r *http.Request) int {
	claims := getClaimsFromRequest(r)
	userID := int(claims["user_id"].(float64))
	return userID
}
