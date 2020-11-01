package test

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/weidonglian/golang-notes-app/app"
	"github.com/weidonglian/golang-notes-app/config"
	"github.com/weidonglian/golang-notes-app/db"
	"github.com/weidonglian/golang-notes-app/logging"
	"github.com/weidonglian/golang-notes-app/model"
	"github.com/weidonglian/golang-notes-app/store"
	"net/http/httptest"

	"github.com/gavv/httpexpect/v2"
	"github.com/onsi/ginkgo"
)

type TestApp struct {
	App    *app.App
	RawAPI *httpexpect.Expect // Raw vanilla request with any header pre-injection
	API    *httpexpect.Expect // With test user's auth header bearer token
	server *httptest.Server
}

func newTestApp(logger *logrus.Logger, cfg config.Config) (*app.App, error) {
	if !config.IsTestMode() {
		panic("NewTestApp should only be used for test application")
	}
	dbSess := db.LoadSessionPool(logger, cfg).ForkNewSession()
	return app.NewAppWith(logger, cfg, dbSess)
}

// Each new test app will fork a new db session and will be cleanup after suite test.
func NewTestAppAndServe() TestApp {
	// Mandatory
	config.SetTestMode()

	logger := logging.NewLogger()
	cfg := config.GetConfig()

	logger.Infof("Creating a new test app")

	if tapp, err := newTestApp(logger, cfg); err != nil {
		panic("Failed to create the test app")
	} else {
		// run server using httptest
		server := httptest.NewServer(tapp.Router())

		// create httpexpect instance
		rawapi := httpexpect.New(ginkgo.GinkgoT(), server.URL)
		testUser := tapp.GetStore().Users.FindByName("test")
		if testUser == nil {
			panic("'test' user has not been injected into the database")
		}
		var testUserToken string
		if token, err := tapp.GetAuth().CreateToken(testUser.ID); err != nil {
			panic("failed to create test user token")
		} else {
			testUserToken = token
		}
		api := rawapi.Builder(func(req *httpexpect.Request) {
			req.WithHeader("Authorization", "Bearer "+testUserToken)
		})

		return TestApp{
			App:    tapp,
			RawAPI: rawapi,
			API:    api,
			server: server,
		}
	}
}

func newReqGraphQLPayload(query string, opts []interface{}) map[string]interface{} {
	var payload map[string]interface{}
	switch len(opts) {
	case 0:
		payload = map[string]interface{}{
			"query": query,
		}
	case 1:
		payload = map[string]interface{}{
			"query":         query,
			"operationName": nil,
			"variables":     opts[0].(map[string]interface{}),
		}
	case 2:
		payload = map[string]interface{}{
			"query":         query,
			"variables":     opts[0].(map[string]interface{}),
			"operationName": opts[1].(string),
		}
	default:
		panic("Only 1 or 2 optional arguments are allowed.")
	}
	return payload
}

func (t *TestApp) GraphPost(query string, opts ...interface{}) *httpexpect.Object {
	payload := newReqGraphQLPayload(query, opts)
	return t.API.POST("/graphql").WithJSON(payload).
		Expect().
		JSON().Object()
}

func (t *TestApp) GraphMustData(query string, opts ...interface{}) *httpexpect.Object {
	payload := newReqGraphQLPayload(query, opts)
	return t.API.POST("/graphql").WithJSON(payload).
		Expect().
		JSON().Object().NotContainsKey("error").ContainsKey("data").Value("data").Object()
}

func (t *TestApp) GraphMustError(query string, opts ...interface{}) *httpexpect.Object {
	payload := newReqGraphQLPayload(query, opts)
	return t.API.POST("/graphql").WithJSON(payload).
		Expect().
		JSON().Object().NotContainsKey("data").ContainsKey("error").Value("error").Object()
}

func (t *TestApp) Close() {
	t.server.Close()
	t.App.Close()
}

func FillDataToStore(s *store.Store, userName string, notes []model.NoteWithTodos) {
	user := s.Users.FindByName(userName)
	if user == nil {
		panic(fmt.Sprintf("could not find user:%s", userName))
	}
	userID := user.ID
	notesStore := s.Notes
	todosStore := s.Todos
	for i := range notes {
		notes[i].UserID = userID
		// note
		if createdNote, err := notesStore.Create(*notes[i].Note); err != nil {
			panic(fmt.Sprintf("failed to create note: %v with error:%v", notes[i], err))
		} else {
			notes[i].ID = createdNote.ID
		}
		// todos
		for j := range notes[i].Todos {
			notes[i].Todos[j].NoteID = notes[i].ID
			if createdTodo, err := todosStore.Create(notes[i].Todos[j]); err != nil {
				panic(fmt.Sprintf("failed to create todo: %d and %d", i, j))
			} else {
				notes[i].Todos[j].ID = createdTodo.ID
			}
		}
	}
}
