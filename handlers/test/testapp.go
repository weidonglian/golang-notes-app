package test

import (
	"fmt"
	"github.com/weidonglian/golang-notes-app/app"
	"github.com/weidonglian/golang-notes-app/config"
	"github.com/weidonglian/golang-notes-app/logging"
	"github.com/weidonglian/golang-notes-app/model"
	"github.com/weidonglian/golang-notes-app/store"
	"net/http/httptest"

	"github.com/gavv/httpexpect/v2"
	"github.com/onsi/ginkgo"
)

type HandlerTestApp struct {
	App    *app.App
	RawAPI *httpexpect.Expect // Raw vanilla request with any header pre-injection
	API    *httpexpect.Expect // With test user's auth header bearer token
	server *httptest.Server
}

// Each new test app will fork a new db session and will be cleanup after suite test.
func NewTestAppAndServe() HandlerTestApp {
	// Mandatory
	config.SetTestMode()

	logger := logging.NewLogger()
	cfg := config.GetConfig()

	logger.Infof("Creating a new test app")

	if tapp, err := app.NewTestApp(logger, cfg); err != nil {
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

		return HandlerTestApp{
			App:    tapp,
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
