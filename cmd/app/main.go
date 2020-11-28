package main

import (
	"github.com/weidonglian/notes-app/config"
	"github.com/weidonglian/notes-app/internal/app"
	"github.com/weidonglian/notes-app/pkg/logging"
	"log"
)

func main() {

	logger := logging.NewLogger()
	cfg := config.GetConfig()

	logger.Infof("The app is running in '%s' mode", config.GetAppMode())

	app, err := app.NewApp(logger, cfg)

	if err != nil {
		logger.Errorf("Failed to create the main app: %s", err.Error())
		return
	}

	if err := app.Serve(); err != nil {
		log.Printf("failed to serve:+%v\n", err)
	}
}
