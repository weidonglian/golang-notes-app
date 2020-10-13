package store_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/weidonglian/golang-notes-app/db"
	"github.com/weidonglian/golang-notes-app/model"
	"github.com/weidonglian/golang-notes-app/store"
)

var _ = Describe("Store", func() {
	var (
		dbSession *db.Session
		sto       *store.Store
	)

	BeforeEach(func() {
		dbSession = dbSessionPool.ForkNewSession()
		if s, err := store.NewStore(dbSession); err != nil {
			panic(err)
		} else {
			sto = s
		}
	})

	AfterEach(func() {
		if err := dbSession.Close(); err != nil {
			panic(err)
		}
	})

	Describe("UsersStore", func() {
		var (
			users = []model.User{
				{
					Username: "u1",
					Password: "p1",
					Role:     model.UserRoleUser,
				},
				{
					Username: "u2",
					Password: "p2",
					Role:     model.UserRoleUser,
				},
				{
					Username: "u3",
					Password: "p3",
					Role:     model.UserRoleUser,
				},
			}
		)

		It("Create, Find and DeleteAll", func() {
			usersStore := sto.Users
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
			err := usersStore.DeleteAll()
			Expect(err).ToNot(HaveOccurred())
			for _, user := range users {
				foundUser := usersStore.FindByName(user.Username)
				Expect(foundUser).To(BeNil())
			}
		})

		It("Find and Delete", func() {
			usersStore := sto.Users
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
				err := usersStore.Delete(foundUser.ID)
				Expect(err).NotTo(HaveOccurred())
				foundUser = usersStore.FindByName(user.Username)
				Expect(foundUser).To(BeNil())
			}
		})
	})

	Describe("NotesStore", func() {
		var (
			testUserId int
			notes      []model.Note
		)

		BeforeEach(func() {
			testUserId = sto.Users.FindByName("test").ID
			notes = []model.Note{
				{
					Name:   "n1",
					UserID: testUserId,
				},
				{
					Name:   "n2",
					UserID: testUserId,
				},
				{
					Name:   "n3",
					UserID: testUserId,
				},
			}
		})

		It("Create, Find and DeleteAll", func() {
			notesStore := sto.Notes
			By("Creat non-existent notes should be pleasant")
			for _, note := range notes {
				noteId, err := notesStore.Create(note)
				Expect(err).NotTo(HaveOccurred())
				foundNote := notesStore.FindByID(noteId)
				Expect(foundNote.Name).To(Equal(note.Name))
				Expect(foundNote.ID).To(Equal(noteId))
				Expect(foundNote.UserID).To(Equal(note.UserID))
			}

			By("Note name is unique should allow duplications")
			for _, note := range notes {
				_, err := notesStore.Create(note)
				Expect(err).NotTo(HaveOccurred())
			}

			By("Should be possible to clear all the notes")
			err := notesStore.DeleteAll()
			Expect(err).ToNot(HaveOccurred())
			for _, note := range notes {
				foundNote := notesStore.FindByName(note.Name)
				Expect(foundNote).To(BeNil())
			}
		})

		It("Find and Delete", func() {
			notesStore := sto.Notes
			By("Creat notes should be found correctly by id")
			for _, note := range notes {
				noteId, err := notesStore.Create(note)
				Expect(err).NotTo(HaveOccurred())
				foundNote := notesStore.FindByID(noteId)
				Expect(foundNote).NotTo(BeNil())
				Expect(foundNote.Name).To(Equal(note.Name))
			}

			By("Find by name and remove by id")
			for _, note := range notes {
				foundNotes := notesStore.FindByName(note.Name)
				Expect(len(foundNotes) > 0).To(BeTrue())
			}
		})
	})
})
