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

	Describe("UsersStore", func() {
		var (
			users = []model.User{
				model.User{
					Username: "u1",
					Password: "p1",
					Role:     store.UserRoleUser,
				},
				model.User{
					Username: "u2",
					Password: "p2",
					Role:     store.UserRoleUser,
				},
				model.User{
					Username: "u3",
					Password: "p3",
					Role:     store.UserRoleUser,
				},
			}
		)

		It("Create, Find and Clear", func() {
			By("Creat non-existent users should be pleasant")
			for _, user := range users {
				userId, err := usersStore.Create(user)
				Expect(err).NotTo(HaveOccurred())
				foundUser := usersStore.FindByID(userId)
				Expect(foundUser.Username).To(Equal(user.Username))
				Expect(store.CheckPassword(foundUser.Password, user.Password)).To(BeTrue())
				Expect(foundUser.Role).To(Equal(user.Role))
			}

			By("Create existent users should not be allowed")
			for _, user := range users {
				_, err := usersStore.Create(user)
				Expect(err).To(HaveOccurred())
			}

			By("Should be possible to clear all the users")
			err := usersStore.Clear()
			Expect(err).ToNot(HaveOccurred())
			for _, user := range users {
				foundUser := usersStore.FindByName(user.Username)
				Expect(foundUser).To(BeNil())
			}
		})

		It("Find and Remove", func() {
			By("Creat users should be found correctly by id")
			for _, user := range users {
				userId, err := usersStore.Create(user)
				Expect(err).NotTo(HaveOccurred())
				foundUser := usersStore.FindByID(userId)
				Expect(foundUser).NotTo(BeNil())
				Expect(foundUser.Username).To(Equal(user.Username))
			}

			By("Find by name and remove by id")
			for _, user := range users {
				foundUser := usersStore.FindByName(user.Username)
				Expect(foundUser).NotTo(BeNil())
				Expect(foundUser.Username).To(Equal(user.Username))
				foundUserFromId := usersStore.FindByID(foundUser.ID)
				Expect(foundUserFromId).To(Equal(foundUser))
				err := usersStore.Remove(foundUser.ID)
				Expect(err).NotTo(HaveOccurred())
				foundUser = usersStore.FindByName(user.Username)
				Expect(foundUser).To(BeNil())
			}
		})
	})
})
