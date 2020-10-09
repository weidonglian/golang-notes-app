package main

import (
	"log"

	"github.com/weidonglian/golang-notes-app/app"
)

func main() {
	app, err := app.NewApp()
	if err != nil {
		log.Print("Failed to create App")
		return
	}

	app.Serve()
}
