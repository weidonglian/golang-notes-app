package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/weidonglian/notes-app/internal/graphql/gmodel"
	"github.com/weidonglian/notes-app/internal/lib"
	"github.com/weidonglian/notes-app/internal/model"
	"github.com/weidonglian/notes-app/internal/pubsub"
)

func (r *mutationResolver) AddNote(ctx context.Context, input gmodel.AddNoteInput) (*gmodel.Note, error) {
	if input.Name == "" {
		return nil, errors.New(`'name' field can not be empty`)
	}

	n, err := r.store.Notes.Create(model.Note{
		Name:   input.Name,
		UserID: lib.GetUserId(ctx),
	})

	if err != nil {
		return nil, err
	}

	gnote := NewGNote(n, make([]model.Todo, 0))
	r.publisher.Publish(ctx, pubsub.EventNoteCreate, gnote)
	return gnote, nil
}

func (r *mutationResolver) UpdateNote(ctx context.Context, input gmodel.UpdateNoteInput) (*gmodel.Note, error) {
	if input.Name == "" {
		return nil, errors.New("'name' field can not be empty")
	}
	n, err := r.store.Notes.Update(input.ID, input.Name, lib.GetUserId(ctx))

	if err != nil {
		return nil, err
	}
	gnote := NewGNote(n, r.store.Todos.FindByNoteID(n.ID))
	r.publisher.Publish(ctx, pubsub.EventNoteUpdate, gnote)
	return gnote, nil
}

func (r *mutationResolver) DeleteNote(ctx context.Context, input *gmodel.DeleteNoteInput) (*gmodel.DeleteNotePayload, error) {
	id, err := r.store.Notes.Delete(input.ID, lib.GetUserId(ctx))
	if err != nil {
		return nil, fmt.Errorf("unprocessable entity with 'id' %d", input.ID)
	}

	payload := &gmodel.DeleteNotePayload{
		ID: id,
	}
	r.publisher.Publish(ctx, pubsub.EventNoteDelete, payload)
	return payload, nil
}

func (r *queryResolver) Notes(ctx context.Context) ([]*gmodel.Note, error) {
	notes := r.store.Notes.FindByUserID(lib.GetUserId(ctx))
	gnotes := make([]*gmodel.Note, len(notes))
	for i := range notes {
		gnotes[i] = NewGNote(&notes[i], r.store.Todos.FindByNoteID(notes[i].ID))
	}
	return gnotes, nil
}

func (r *queryResolver) Note(ctx context.Context, id int) (*gmodel.Note, error) {
	n := r.store.Notes.FindByID(id, lib.GetUserId(ctx))
	if n != nil {
		return NewGNote(n, r.store.Todos.FindByNoteID(n.ID)), nil
	} else {
		return nil, fmt.Errorf("failed to find a note with id '%d'", id)
	}
}

func (r *subscriptionResolver) NoteAdded(ctx context.Context) (<-chan *gmodel.Note, error) {
	chanNote := make(chan *gmodel.Note, 1)

	err := r.subscriber.Subscribe(ctx, pubsub.EventNoteCreate, func(msg *nats.Msg) {
		var gnote gmodel.Note
		if err := json.Unmarshal(msg.Data, &gnote); err != nil {
			r.logger.Errorf("unable to unmarshal pubsub event data: %s with type %T", msg.Subject, gnote)
		}
		chanNote <- &gnote
	})
	if err != nil {
		return nil, err
	}

	return chanNote, nil
}

func (r *subscriptionResolver) NoteUpdated(ctx context.Context) (<-chan *gmodel.Note, error) {
	chanNote := make(chan *gmodel.Note, 1)

	err := r.subscriber.Subscribe(ctx, pubsub.EventNoteUpdate, func(msg *nats.Msg) {
		var gnote gmodel.Note
		if err := json.Unmarshal(msg.Data, &gnote); err != nil {
			r.logger.Errorf("unable to unmarshal pubsub event data: %s type %T", msg.Subject, gnote)
		}
		chanNote <- &gnote
	})
	if err != nil {
		return nil, err
	}

	return chanNote, nil
}

func (r *subscriptionResolver) NoteDeleted(ctx context.Context) (<-chan *gmodel.DeleteNotePayload, error) {
	chanDeletePayload := make(chan *gmodel.DeleteNotePayload, 1)

	err := r.subscriber.Subscribe(ctx, pubsub.EventNoteDelete, func(msg *nats.Msg) {
		var payload gmodel.DeleteNotePayload
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
