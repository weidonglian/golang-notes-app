package handlers_test

import (
	"github.com/weidonglian/golang-notes-app/app"
	"github.com/weidonglian/golang-notes-app/config"
	"github.com/weidonglian/golang-notes-app/logging"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect/v2"
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

var testApp HandlerTestApp

var _ = BeforeSuite(func() {
	testApp = NewTestAppAndServe(GinkgoT())
})

var _ = AfterSuite(func() {
	testApp.Close()
})

type HandlerTestApp struct {
	App    *app.App
	API    *httpexpect.Expect
	server *httptest.Server
}

func NewTestAppAndServe(t GinkgoTInterface) HandlerTestApp {
	if !config.IsTestMode() {
		panic("NewTestApp should only be allowed in test mode.")
	}

	logger := logging.NewLogger()
	cfg := config.GetConfig()

	logger.Infof("Creating a new test app")

	if app, err := app.NewApp(logger, cfg); err != nil {
		panic("Failed to create the test app")
	} else {
		// run server using httptest
		server := httptest.NewServer(app.Router())

		// create httpexpect instance
		api := httpexpect.New(t, server.URL)

		return HandlerTestApp{
			App:    app,
			API:    api,
			server: server,
		}
	}
}

func (ta HandlerTestApp) Close() {
	ta.server.Close()
	ta.App.Close()
}
