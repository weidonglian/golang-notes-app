package util

import (
	"github.com/weidonglian/golang-notes-app/auth"
	"net/http"
)

func GetUserIDFromRequest(r *http.Request) int {
	claims := auth.GetClaimsFromRequest(r.Context())
	userID := int(claims["user_id"].(float64))
	return userID
}
