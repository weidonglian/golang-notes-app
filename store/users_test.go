package store_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/weidonglian/golang-notes-app/config"
	"github.com/weidonglian/golang-notes-app/db"
	"github.com/weidonglian/golang-notes-app/logging"
	"github.com/weidonglian/golang-notes-app/store"
)

func init() {
	config.SetTestMode()
}

var _ = Describe("Users", func() {
	var (
		usersStore   store.UsersStore
		storeContext *store.StoreContext
	)

	BeforeEach(func() {
		sess, err := db.NewSession(logging.NewLogger(), config.GetConfig())
		if err != nil {
			panic(err)
		}

		storeContext = &store.StoreContext{
			sess,
		}

		usersStore = store.NewUsersStore(storeContext)
	})

	AfterEach(func() {
		if storeContext != nil {
			storeContext.Session.Close()
		}
	})

	Describe("USERS Create", func() {
		Expect(usersStore).To(Equal(nil))
	})
})
