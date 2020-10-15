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
	RawAPI *httpexpect.Expect
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
		rawapi := httpexpect.New(GinkgoT(), server.URL)
		testUser := app.GetStore().Users.FindByName("test")
		if testUser == nil {
			panic("'test' user has not been injected into the database")
		}
		var testUserToken string
		if token, err := app.GetAuth().CreateToken(testUser.ID); err != nil {
			panic("failed to create test user token")
		} else {
			testUserToken = token
		}
		api := rawapi.Builder(func(req *httpexpect.Request) {
			req.WithHeader("Authorization", "Bearer "+testUserToken)
		})

		return HandlerTestApp{
			App:    app,
			RawAPI: rawapi,
			API:    api,
			server: server,
		}
	}
}

func (ta HandlerTestApp) Close() {
	ta.server.Close()
	ta.App.Close()
}
