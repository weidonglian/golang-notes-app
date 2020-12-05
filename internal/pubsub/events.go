package pubsub

const (
	// when no message should be published
	NoEvent SubjectKey = ""

	// Notes
	NoteCreateEvent SubjectKey = "app.entity.note.create"
	NoteUpdateEvent SubjectKey = "app.entity.note.update"
	NoteDeleteEvent SubjectKey = "app.entity.note.delete"

	// Todos
	TodoCreateEvent SubjectKey = "app.entity.todo.create"
	TodoUpdateEvent SubjectKey = "app.entity.todo.update"
	TodoDeleteEvent SubjectKey = "app.entity.todo.delete"
)
