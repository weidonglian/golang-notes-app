package lib

import (
	"context"
	"net/http"
)

func GetUserIDFromRequest(r *http.Request) int {
	claims := GetClaimsFromRequest(r.Context())
	userID := int(claims["user_id"].(float64))
	return userID
}

func GetUserId(ctx context.Context) int {
	claims := GetClaimsFromRequest(ctx)
	userID := int(claims["user_id"].(float64))
	return userID
}
