package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/weidonglian/golang-notes-app/graph/gmodel"
)

func (r *mutationResolver) AddNote(ctx context.Context, input gmodel.AddNoteInput) (*gmodel.Note, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateNote(ctx context.Context, input gmodel.UpdateNoteInput) (*gmodel.Note, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeleteNote(ctx context.Context, input *gmodel.DeleteNoteInput) (*gmodel.Note, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Notes(ctx context.Context) ([]*gmodel.Note, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Note(ctx context.Context, id string) (*gmodel.Note, error) {
	panic(fmt.Errorf("not implemented"))
}
