package test

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/weidonglian/notes-app/config"
	"github.com/weidonglian/notes-app/internal/app"
	"github.com/weidonglian/notes-app/internal/db"
	"github.com/weidonglian/notes-app/internal/model"
	"github.com/weidonglian/notes-app/internal/store"
	"github.com/weidonglian/notes-app/pkg/logging"
	"net/http/httptest"

	"github.com/gavv/httpexpect/v2"
	"github.com/onsi/ginkgo"
)

type MockApp struct {
	App    *app.App
	RawAPI *httpexpect.Expect // Raw vanilla request with any header pre-injection
	API    *httpexpect.Expect // With test user's auth header bearer token
	server *httptest.Server
}

func newTestApp(logger *logrus.Logger, cfg config.Config) (*app.App, error) {
	if !config.IsTestMode() {
		panic("NewTestApp should only be used for test application")
	}
	dbSess := db.NewForkedRandomSession(logger, cfg)
	return app.NewAppWith(logger, cfg, dbSess)
}

// Each new test app will fork a new db session and will be cleanup after suite test.
func NewMockApp() MockApp {
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

		return MockApp{
			App:    tapp,
			RawAPI: rawapi,
			API:    api,
			server: server,
		}
	}
}

// graphql request format, it has to be `POST` method with json body as:
//
// {
//  "query": "...",
//  "operationName": "...",
//  "variables": { "myVariable": "someValue", ... }
// }
func newReqGraphqlPayload(query string, opts []interface{}) map[string]interface{} {
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
			"variables":     opts[0],
		}
	case 2:
		payload = map[string]interface{}{
			"query":         query,
			"variables":     opts[0],
			"operationName": opts[1],
		}
	default:
		panic("Only 1 or 2 optional arguments are allowed.")
	}
	return payload
}

// graphql response always http succeed with json body as:
//
// {
//  "data": { ... },
//  "errors": [ ... ]
// }
// GraphqlPost is a bit raw graphql API that you need to extract the `data` and `errors`
// `opts` only supports two optional arguments, 1st argument is `variables` should be
// a `struct` with json tags or a `map[string]interface{}` that can be converted to
// a json. The 2nd argument must be a string for `operationName`. If no argument is
// provided then it will include no `variables` and `operationName`.
func (t *MockApp) GraphqlPost(query string, opts ...interface{}) *httpexpect.Object {
	payload := newReqGraphqlPayload(query, opts)
	return t.API.POST("/graphql").WithJSON(payload).
		Expect().
		JSON().Object()
}

// sugar that expects the `data` only no error. It will return the `data` json object.
func (t *MockApp) GraphqlMustData(query string, opts ...interface{}) *httpexpect.Object {
	payload := newReqGraphqlPayload(query, opts)
	return t.API.POST("/graphql").WithJSON(payload).
		Expect().
		JSON().Object().NotContainsKey("errors").ContainsKey("data").Value("data").Object()
}

// sugar to expects `errors` but `data` could exist with `null`. It will return the `errors` json array.
func (t *MockApp) GraphqlMustError(query string, opts ...interface{}) *httpexpect.Array {
	payload := newReqGraphqlPayload(query, opts)
	return t.API.POST("/graphql").WithJSON(payload).
		Expect().
		JSON().Object().ContainsKey("errors").Value("errors").Array()
}

func GraphqlWithInput(v interface{}) map[string]interface{} {
	return map[string]interface{}{
		"input": v,
	}
}

func (t *MockApp) Close() {
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
