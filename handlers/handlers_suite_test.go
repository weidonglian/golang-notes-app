package handlers_test

import (
	"github.com/weidonglian/golang-notes-app/app"
	"github.com/weidonglian/golang-notes-app/config"
	"github.com/weidonglian/golang-notes-app/logging"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func init() {
	config.SetTestMode()
}

func TestHandlers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Handlers Suite")
}

type TestApp struct {
	App *app.App
}

func NewTestAppAndServe() TestApp {
	if !config.IsTestMode() {
		panic("NewTestApp should only be allowed in test mode.")
	}

	logger := logging.NewLogger()
	cfg := config.GetConfig()

	logger.Infof("Creating a new test app")

	if app, err := app.NewApp(logger, cfg); err != nil {
		panic("Failed to create the test app")
	} else {
		go app.Serve()
		return TestApp{
			App: app,
		}
	}
}
