package main

import (
	"github.com/weidonglian/notes-app/config"
	"github.com/weidonglian/notes-app/internal/app"
	"github.com/weidonglian/notes-app/internal/lib"
	"log"
)

func main() {

	logger := lib.NewLogger()
	cfg := config.DefaultConfig()
	config.MustParseFromEnv(cfg)

	logger.Info("The app is running now")

	app, err := app.NewApp(logger, *cfg)

	if err != nil {
		logger.Errorf("Failed to create the main app: %s", err.Error())
		return
	}

	if err := app.Serve(); err != nil {
		log.Printf("failed to serve:+%v\n", err)
	}
}
