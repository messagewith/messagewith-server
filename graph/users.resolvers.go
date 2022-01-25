package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"messagewith-server/graph/generated"
	"messagewith-server/graph/model"
	"messagewith-server/users"
)

func (r *mutationResolver) CreateUser(ctx context.Context, userInput model.UserInput) (*model.User, error) {
	return users.Service.CreateUser(ctx, &userInput)
}

func (r *mutationResolver) ChangeUserPassword(ctx context.Context, email string, token string, newPassword string) (*model.User, error) {
	return users.Service.ChangePassword(ctx, email, token, newPassword)
}

func (r *mutationResolver) GenerateChangeUserPasswordToken(ctx context.Context, email string) (*string, error) {
	return users.Service.GenerateChangePasswordToken(ctx, email)
}

func (r *queryResolver) Users(ctx context.Context, filter *model.UserFilter) ([]*model.User, error) {
	return users.Service.GetUsers(ctx, filter)
}

func (r *queryResolver) User(ctx context.Context, id *string, email *string) (*model.User, error) {
	return users.Service.GetUser(ctx, id, email, nil)
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
