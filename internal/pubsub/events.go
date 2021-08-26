package pubsub

const (
	// when no message should be published
	NoEvent SubjectKey = ""

	// Notes
	EventNoteCreate SubjectKey = "app.entity.note.create"
	EventNoteUpdate SubjectKey = "app.entity.note.update"
	EventNoteDelete SubjectKey = "app.entity.note.delete"

	// Todos
	EventTodoCreate SubjectKey = "app.entity.todo.create"
	EventTodoUpdate SubjectKey = "app.entity.todo.update"
	EventTodoDelete SubjectKey = "app.entity.todo.delete"
)
