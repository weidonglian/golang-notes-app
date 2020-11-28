package util

import (
	"context"
	"github.com/weidonglian/notes-app/internal/auth"
	"log"
	"net/http"
	"os"
	"os/signal"
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

func NewShutdownContext() context.Context {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		oscall := <-c
		log.Printf("system call:%+v", oscall)
		cancel()
	}()

	return ctx
}
