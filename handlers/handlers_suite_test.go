package handlers_test

import (
	"github.com/weidonglian/golang-notes-app/app"
	"github.com/weidonglian/golang-notes-app/config"
	"github.com/weidonglian/golang-notes-app/db"
	"github.com/weidonglian/golang-notes-app/logging"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect/v2"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestHandlers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Handlers Suite")
}

var _ = AfterSuite(func() {
	db.UnloadSessionPool()
})

type HandlerTestApp struct {
	App    *app.App
	API    *httpexpect.Expect
	server *httptest.Server
}

// Each new test app will fork a new db session and will be cleanup after suite test.
func NewTestAppAndServe() HandlerTestApp {
	// Mandatory
	config.SetTestMode()

	logger := logging.NewLogger()
	cfg := config.GetConfig()

	logger.Infof("Creating a new test app")

	if app, err := app.NewTestApp(logger, cfg); err != nil {
		panic("Failed to create the test app")
	} else {
		// run server using httptest
		server := httptest.NewServer(app.Router())

		// create httpexpect instance
		api := httpexpect.New(GinkgoT(), server.URL)

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
