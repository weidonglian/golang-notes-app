package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"fmt"

	"github.com/weidonglian/notes-app/internal/graphql/gmodel"
	"github.com/weidonglian/notes-app/internal/model"
	"github.com/weidonglian/notes-app/pkg/util"
)

func (r *mutationResolver) AddNote(ctx context.Context, input gmodel.AddNoteInput) (*gmodel.Note, error) {
	if input.Name == "" {
		return nil, errors.New(`'name' field can not be empty`)
	}

	n, err := r.store.Notes.Create(model.Note{
		Name:   input.Name,
		UserID: util.GetUserId(ctx),
	})

	if err != nil {
		return nil, err
	}

	return NewGNote(n, make([]model.Todo, 0)), nil
}

func (r *mutationResolver) UpdateNote(ctx context.Context, input gmodel.UpdateNoteInput) (*gmodel.Note, error) {
	if input.Name == "" {
		return nil, errors.New("'name' field can not be empty")
	}
	n, err := r.store.Notes.Update(input.ID, input.Name, util.GetUserId(ctx))

	if err != nil {
		return nil, err
	}

	return NewGNote(n, r.store.Todos.FindByNoteID(n.ID)), nil
}

func (r *mutationResolver) DeleteNote(ctx context.Context, input *gmodel.DeleteNoteInput) (*gmodel.DeleteNotePayload, error) {
	id, err := r.store.Notes.Delete(input.ID, util.GetUserId(ctx))
	if err != nil {
		return nil, fmt.Errorf("unprocessable entity with 'id' %d", input.ID)
	}

	return &gmodel.DeleteNotePayload{
		ID: id,
	}, nil
}

func (r *queryResolver) Notes(ctx context.Context) ([]*gmodel.Note, error) {
	notes := r.store.Notes.FindByUserID(util.GetUserId(ctx))
	gnotes := make([]*gmodel.Note, len(notes))
	for i := range notes {
		gnotes[i] = NewGNote(&notes[i], r.store.Todos.FindByNoteID(notes[i].ID))
	}
	return gnotes, nil
}

func (r *queryResolver) Note(ctx context.Context, id int) (*gmodel.Note, error) {
	n := r.store.Notes.FindByID(id, util.GetUserId(ctx))
	if n != nil {
		return NewGNote(n, r.store.Todos.FindByNoteID(n.ID)), nil
	} else {
		return nil, fmt.Errorf("failed to find a note with id '%d'", id)
	}
}
