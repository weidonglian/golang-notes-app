package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/weidonglian/golang-notes-app/graph/model"
)

func (r *mutationResolver) AddNote(ctx context.Context, input model.AddNoteInput) (*model.Note, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateNote(ctx context.Context, input model.UpdateNoteInput) (*model.Note, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeleteNote(ctx context.Context, input *model.DeleteNoteInput) (*model.Note, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Notes(ctx context.Context) ([]*model.Note, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Note(ctx context.Context, id string) (*model.Note, error) {
	panic(fmt.Errorf("not implemented"))
}
