package main

import (
	"github.com/weidonglian/golang-notes-app/app"
	"github.com/weidonglian/golang-notes-app/logging"
)

func main() {

	logger := logging.NewLogger()

	app, err := app.NewApp(logger)

	if err != nil {
		logger.Error("Failed to create the main app")
		return
	}

	app.Serve()
}
