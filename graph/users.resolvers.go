package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"awesomeProject/graph/generated"
	"awesomeProject/graph/model"
	"context"
)

func (r *mutationResolver) CreateUser(ctx context.Context, userInput model.UserInput) (*model.User, error) {
	return r.usersService.CreateUser(ctx, &userInput)
}

func (r *queryResolver) Users(ctx context.Context, filter *model.UserFilter) ([]*model.User, error) {
	return r.usersService.GetUsers(ctx, filter)
}

func (r *queryResolver) User(ctx context.Context, id *string, email *string) (*model.User, error) {
	return r.usersService.GetUser(ctx, id, email)
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
