package util

import (
	"context"
	"github.com/weidonglian/notes-app/internal/auth"
	"net/http"
)

func GetUserIDFromRequest(r *http.Request) int {
	claims := auth.GetClaimsFromRequest(r.Context())
	userID := int(claims["user_id"].(float64))
	return userID
}

func GetUserId(ctx context.Context) int {
	claims := auth.GetClaimsFromRequest(ctx)
	userID := int(claims["user_id"].(float64))
	return userID
}
