package store_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
	"github.com/weidonglian/notes-app/config"
	"github.com/weidonglian/notes-app/internal/db"
	"github.com/weidonglian/notes-app/internal/model"
	"github.com/weidonglian/notes-app/internal/store"
	"github.com/weidonglian/notes-app/internal/test"
	"github.com/weidonglian/notes-app/pkg/logging"
	"github.com/weidonglian/notes-app/pkg/util"
)

var _ = Describe("Store", func() {
	var (
		dbSession db.Session
		sto       *store.Store
	)

	// we need to set up a new db session for different stores
	BeforeEach(func() {
		logger := logging.NewLogger()
		logger.SetLevel(logrus.WarnLevel)
		dbSession = db.NewForkedSession(logger, *config.DefaultTestConfig())
		if s, err := store.NewStore(dbSession, logger); err != nil {
			panic(err)
		} else {
			sto = s
		}
		test.LoadTestUsers(sto)
	})

	AfterEach(func() {
		if err := dbSession.Close(); err != nil {
			panic(err)
		}
	})

	// UsesStore tests
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
				createdUser, err := usersStore.Create(user)
				Expect(err).NotTo(HaveOccurred())
				foundUser := usersStore.FindByID(createdUser.ID)
				Expect(foundUser.Username).To(Equal(user.Username))
				Expect(util.CheckPassword(foundUser.Password, user.Password)).To(BeTrue())
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
				createdUser, err := usersStore.Create(user)
				Expect(err).NotTo(HaveOccurred())
				foundUser := usersStore.FindByID(createdUser.ID)
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

	// NotesStore
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
				createdNote, err := notesStore.Create(note)
				Expect(err).NotTo(HaveOccurred())
				foundNote := notesStore.FindByID(createdNote.ID, testUserId)
				Expect(foundNote).NotTo(BeNil())
				Expect(foundNote.Name).To(Equal(note.Name))
				Expect(foundNote.ID).To(Equal(createdNote.ID))
				Expect(foundNote.UserID).To(Equal(note.UserID))
			}

			By("Note name is unique should allow duplications")
			for _, note := range notes {
				_, err := notesStore.Create(note)
				Expect(err).NotTo(HaveOccurred())
			}

			By("Should be possible to clear all the notes")
			err := notesStore.DeleteAll(testUserId)
			Expect(err).ToNot(HaveOccurred())
			for _, note := range notes {
				foundNote := notesStore.FindByName(note.Name, testUserId)
				Expect(foundNote).To(BeEmpty())
			}
		})

		It("Find and Delete", func() {
			notesStore := sto.Notes
			By("Creat notes should be found correctly by id")
			for _, note := range notes {
				createdNote, err := notesStore.Create(note)
				Expect(err).NotTo(HaveOccurred())
				foundNote := notesStore.FindByID(createdNote.ID, testUserId)
				Expect(foundNote).NotTo(BeNil())
				Expect(foundNote.Name).To(Equal(note.Name))
			}

			By("Find by name and remove by id")
			for _, note := range notes {
				foundNotes := notesStore.FindByName(note.Name, testUserId)
				Expect(len(foundNotes) > 0).To(BeTrue())
			}
		})
	})

	// TodosStore
	Describe("TodosStore", func() {
		var (
			testUserId int
			notes      []model.Note
			todos      []model.Todo
		)

		BeforeEach(func() {
			// create notes for current test user
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
			// add notes to db
			notesStore := sto.Notes
			for i := range notes {
				n := &notes[i]
				if createdNote, err := notesStore.Create(*n); err != nil {
					panic(err)
				} else {
					n.ID = createdNote.ID
				}
			}
			// init todo list
			todos = []model.Todo{
				{
					Name: "t1",
					Done: false,
				},
				{
					Name: "t1",
					Done: true,
				},
				{
					Name: "t3",
					Done: false,
				},
				{
					Name: "t4",
					Done: true,
				},
			}
		})

		It("Create, Find and DeleteAll", func() {
			todosStore := sto.Todos
			By("Creat todoNames in an existing note should always be pleasant")
			for _, note := range notes {
				for _, todo := range todos {
					By("Create a todo into given note should always work")
					todo.NoteID = note.ID
					createdTodo, err := todosStore.Create(todo)
					Expect(err).NotTo(HaveOccurred())

					By("Found by todo id should always work")
					foundTodo := todosStore.FindByID(createdTodo.ID)
					Expect(foundTodo.Name).To(Equal(todo.Name)) // should be exact since search by unique id
					Expect(foundTodo.ID).To(Equal(createdTodo.ID))
					Expect(foundTodo.Done).To(Equal(todo.Done))
				}
			}
		})

		It("Find and Delete", func() {
			todosStore := sto.Todos
			By("Creat todoNames in an existing note should always be pleasant")
			for _, note := range notes {
				for _, todo := range todos {
					todo.NoteID = note.ID
					_, err := todosStore.Create(todo)
					Expect(err).NotTo(HaveOccurred())
				}
			}

			By("Find by name 't1' should have 6 todos")
			Expect(len(todosStore.FindByName("t1"))).To(BeEquivalentTo(6))

			By("Find by name 't2, t3' should have 3 todos")
			Expect(len(todosStore.FindByName("t3"))).To(BeEquivalentTo(3))
			Expect(len(todosStore.FindByName("t4"))).To(BeEquivalentTo(3))

			By("Find by name 'non-existent' should have 0 todos")
			Expect(len(todosStore.FindByName("t2"))).To(BeZero())
			Expect(len(todosStore.FindByName("invalid_xxx"))).To(BeZero())

			By("Find by note id should have 4 todos per id")
			for _, note := range notes {
				foundTodos := todosStore.FindByNoteID(note.ID)
				Expect(len(foundTodos)).To(BeIdenticalTo(4))
				By("Remove all those one by one and should drop to zero")
				for _, todo := range foundTodos {
					todoID, err := todosStore.Delete(todo.ID, note.ID)
					Expect(err).NotTo(HaveOccurred())
					Expect(todoID).To(Equal(todo.ID))
				}
				Expect(len(todosStore.FindByNoteID(note.ID))).To(BeZero())
			}
		})
	})
})
