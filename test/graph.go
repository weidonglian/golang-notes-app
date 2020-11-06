package test

// GraphQL `note`
var QueryNotes = `
query {
  notes {
	id
	name    
	todos {
	  id
	  name
	  done
	  noteId
	}
  }
}
`

var QueryNote = `
query($id: Int!){
  note(id: $id) {
	id
	name    
	todos {
	  id
	  name
	  done
	  noteId
	}
  }
}
`

var MutationAddNote = `
mutation ($input: AddNoteInput!) {
  addNote(input: $input) {
    id
    name
    todos {
	  id
	  name
	  done
	  noteId
    }
  }
}
`

var MutationUpdateNote = `
mutation ($input: UpdateNoteInput!) {
  updateNote(input: $input) {
    id
    name
    todos {
      id
      name
	  done
	  noteId
    }
  }
}
`

var MutationDeleteNote = `mutation ($input: DeleteNoteInput!) {
  deleteNote(input: $input) {
    id
  }
}`

// GraphQL `todo`
var QueryTodos = `
query ($noteId: Int!){
  todos(noteId: $noteId) {
    id
    name
    done
    noteId
  }
}
`

var MutationAddTodo = `
mutation ($input: AddTodoInput!){
  addTodo(input: $input) {
    id
    name
	done
	noteId
  }
}
`

var MutationUpdateTodo = `
mutation ($input: UpdateTodoInput!){
  updateTodo(input: $input) {
    id
    name
    done
    noteId
  }
}
`

var MutationToggleTodo = `
mutation ($input: ToggleTodoInput!){
  toggleTodo(input: $input) {    
    id
	name
	done
	noteId    
  }
}
`

var MutationDeleteTodo = `
mutation ($input: DeleteTodoInput!){
  deleteTodo(input: $input) {
    id
    noteId
  }
}
`
