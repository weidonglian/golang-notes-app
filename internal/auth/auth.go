package auth

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/jwtauth"
	"github.com/weidonglian/notes-app/config"
)

type Auth struct {
	tokenAuth *jwtauth.JWTAuth
}

func (auth Auth) CreateToken(userID int) (string, error) {
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = userID
	atClaims["exp"] = time.Now().Add(time.Hour * 120).Unix()
	return auth.generateToken(atClaims)
}

func (auth Auth) generateToken(claims jwt.MapClaims) (string, error) {
	// For debugging/example purposes, we generate and print
	// a sample jwt token with claims `user_id:123` here:
	_, tokenString, err := auth.tokenAuth.Encode(claims)
	return tokenString, err
}

func (auth Auth) Verifier() func(http.Handler) http.Handler {
	return jwtauth.Verifier(auth.tokenAuth)
}

func (auth Auth) Authenticator() func(http.Handler) http.Handler {
	return jwtauth.Authenticator
}

func (auth Auth) VerifyToken(tokenStr string) (*jwt.Token, error) {
	return jwtauth.VerifyRequest(auth.tokenAuth, nil, func(r *http.Request) string {
		return tokenStr
	})
}

func NewAuth(cfg config.Config) *Auth {
	tokenAuth := jwtauth.New("HS256", []byte(cfg.JWTSecret), nil)
	return &Auth{tokenAuth}
}
