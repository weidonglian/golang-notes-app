package store_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/weidonglian/golang-notes-app/config"
	"github.com/weidonglian/golang-notes-app/db"
	"github.com/weidonglian/golang-notes-app/model"
	"github.com/weidonglian/golang-notes-app/store"
)

func init() {
	config.SetTestMode()
}

var _ = Describe("Users", func() {
	var (
		usersStore store.UsersStore
		dbSession  *db.Session
	)

	BeforeEach(func() {
		dbSession = dbSessionPool.ForkNewSession()
		usersStore = store.NewUsersStore(&store.StoreContext{
			Session: dbSession,
		})
	})

	AfterEach(func() {
		dbSession.Close()
	})

	Describe("USERS Create", func() {
		It("Create user", func() {
			user := model.User{
				Username: "u1",
				Password: "p1",
				Role:     store.UserRoleUser,
			}
			userId, err := usersStore.Create(user)
			Expect(err).NotTo(HaveOccurred())

			foundUser := usersStore.FindByID(userId)
			Expect(foundUser.Username).To(Equal(user.Username))
			Expect(store.CheckPassword(foundUser.Password, user.Password)).To(BeTrue())
			Expect(foundUser.Role).To(Equal(user.Role))
		})
	})
})
