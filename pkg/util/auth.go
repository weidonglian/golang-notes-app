package util

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/jwtauth"
	"golang.org/x/crypto/bcrypt"
)

func GetClaimsFromRequest(ctx context.Context) jwt.MapClaims {
	if _, claims, err := jwtauth.FromContext(ctx); err != nil {
		panic(err)
	} else {
		return claims
	}
}

func HashPassword(pwd string) (string, error) {

	// Use GenerateFromPassword to hash & salt pwd
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash), nil
}

func CheckPassword(hashedPwd string, plainPwd string) bool {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(plainPwd))
	if err != nil {
		return false
	}
	return true
}
