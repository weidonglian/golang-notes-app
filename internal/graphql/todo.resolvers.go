package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"encoding/json"
	"github.com/nats-io/nats.go"

	"github.com/weidonglian/notes-app/internal/graphql/gmodel"
	"github.com/weidonglian/notes-app/internal/lib"
	"github.com/weidonglian/notes-app/internal/model"
	"github.com/weidonglian/notes-app/internal/pubsub"
)

func (r *mutationResolver) AddTodo(ctx context.Context, input gmodel.AddTodoInput) (*gmodel.Todo, error) {
	done := false
	if input.Done != nil {
		done = *input.Done
	}

	todo, err := r.store.Todos.Create(model.Todo{
		Name:   input.Name,
		Done:   done,
		NoteID: input.NoteID,
	})

	if err != nil {
		return nil, lib.ErrorUnprocessableEntity
	}

	gtodo := NewGTodo(todo)
	r.publisher.Publish(ctx, pubsub.EventTodoCreate, gtodo)
	return gtodo, nil
}

func (r *mutationResolver) UpdateTodo(ctx context.Context, input gmodel.UpdateTodoInput) (*gmodel.Todo, error) {
	if r.store.Notes.FindByID(input.NoteID, lib.GetUserId(ctx)) == nil {
		return nil, lib.ErrorUnprocessableEntity
	}

	todo, err := r.store.Todos.Update(input.ID, input.Name, input.Done)

	if err != nil {
		return nil, lib.ErrorUnprocessableEntity
	}

	gtodo := NewGTodo(todo)
	r.publisher.Publish(ctx, pubsub.EventTodoUpdate, gtodo)
	return gtodo, nil
}

func (r *mutationResolver) DeleteTodo(ctx context.Context, input gmodel.DeleteTodoInput) (*gmodel.DeleteTodoPayload, error) {
	if r.store.Notes.FindByID(input.NoteID, lib.GetUserId(ctx)) == nil {
		return nil, lib.ErrorUnprocessableEntity
	}

	id, err := r.store.Todos.Delete(input.ID, input.NoteID)
	if err != nil {
		return nil, lib.ErrorUnprocessableEntity
	}

	payload := &gmodel.DeleteTodoPayload{
		ID:     id,
		NoteID: input.NoteID,
	}
	r.publisher.Publish(ctx, pubsub.EventTodoDelete, payload)
	return payload, nil
}

func (r *mutationResolver) ToggleTodo(ctx context.Context, input gmodel.ToggleTodoInput) (*gmodel.Todo, error) {
	if r.store.Notes.FindByID(input.NoteID, lib.GetUserId(ctx)) == nil {
		return nil, lib.ErrorUnprocessableEntity
	}

	todo, err := r.store.Todos.Toggle(input.ID)

	if err != nil {
		return nil, lib.ErrorUnprocessableEntity
	}

	gtodo := NewGTodo(todo)
	r.publisher.Publish(ctx, pubsub.EventTodoUpdate, gtodo)
	return gtodo, nil
}

func (r *queryResolver) Todos(ctx context.Context, noteID int) ([]*gmodel.Todo, error) {
	if r.store.Notes.FindByID(noteID, lib.GetUserId(ctx)) == nil {
		return nil, lib.ErrorUnprocessableEntity
	}

	todos := r.store.Todos.FindByNoteID(noteID)
	gtodos := make([]*gmodel.Todo, len(todos))
	for i := range todos {
		gtodos[i] = NewGTodo(&todos[i])
	}
	return gtodos, nil
}

func (r *subscriptionResolver) TodoAdded(ctx context.Context) (<-chan *gmodel.Todo, error) {
	chanTodo := make(chan *gmodel.Todo, 1)

	err := r.subscriber.Subscribe(ctx, pubsub.EventTodoCreate, func(msg *nats.Msg) {
		var gtodo gmodel.Todo
		if err := json.Unmarshal(msg.Data, &gtodo); err != nil {
			r.logger.Errorf("unable to unmarshal pubsub event data: %s with type %T", msg.Subject, gtodo)
		}
		chanTodo <- &gtodo
	})
	if err != nil {
		return nil, err
	}

	return chanTodo, nil
}

func (r *subscriptionResolver) TodoUpdated(ctx context.Context) (<-chan *gmodel.Todo, error) {
	chanTodo := make(chan *gmodel.Todo, 1)

	err := r.subscriber.Subscribe(ctx, pubsub.EventTodoUpdate, func(msg *nats.Msg) {
		var gtodo gmodel.Todo
		if err := json.Unmarshal(msg.Data, &gtodo); err != nil {
			r.logger.Errorf("unable to unmarshal pubsub event data: %s with type %T", msg.Subject, gtodo)
		}
		chanTodo <- &gtodo
	})
	if err != nil {
		return nil, err
	}

	return chanTodo, nil
}

func (r *subscriptionResolver) TodoDeleted(ctx context.Context) (<-chan *gmodel.DeleteTodoPayload, error) {
	chanDeletePayload := make(chan *gmodel.DeleteTodoPayload, 1)

	err := r.subscriber.Subscribe(ctx, pubsub.EventTodoDelete, func(msg *nats.Msg) {
		var payload gmodel.DeleteTodoPayload
		if err := json.Unmarshal(msg.Data, &payload); err != nil {
			r.logger.Errorf("unable to unmarshal pubsub event data: %s type %T", msg.Subject, payload)
		}
		chanDeletePayload <- &payload
	})
	if err != nil {
		return nil, err
	}

	return chanDeletePayload, nil
}
