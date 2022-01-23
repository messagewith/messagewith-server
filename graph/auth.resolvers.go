package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"awesomeProject/graph/model"
	"context"
)

func (r *mutationResolver) Logout(ctx context.Context) (*bool, error) {
	return r.authService.Logout(ctx)
}

func (r *mutationResolver) Login(ctx context.Context, email string, password string) (*model.User, error) {
	return r.authService.Login(ctx, email, password)
}

func (r *queryResolver) LoggedUser(ctx context.Context) (*model.User, error) {
	return r.authService.GetLoggedUser(ctx)
}
