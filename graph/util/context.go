package util

import (
	"context"
	"github.com/weidonglian/golang-notes-app/auth"
)

func GetUserId(ctx context.Context) int {
	claims := auth.GetClaimsFromRequest(ctx)
	userID := int(claims["user_id"].(float64))
	return userID
}
