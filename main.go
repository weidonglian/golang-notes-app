package main

import (
	"github.com/weidonglian/notes-app/app"
	"github.com/weidonglian/notes-app/config"
	"github.com/weidonglian/notes-app/logging"
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

	app.Serve()
}
