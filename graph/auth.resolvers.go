package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"messagewith-server/auth"
	errors "messagewith-server/errors"
	"messagewith-server/graph/model"
	"messagewith-server/users"
)

func (r *mutationResolver) Logout(ctx context.Context) (*bool, error) {
	authService := auth.GetService(ctx)
	if user := users.UserFromContext(ctx); user == nil {
		return nil, errors.ErrUserNotLoggedIn
	}

	return authService.Logout(ctx)
}

func (r *mutationResolver) Login(ctx context.Context, email string, password string) (*model.User, error) {
	authService := auth.GetService(ctx)
	if user := users.UserFromContext(ctx); user != nil {
		return nil, errors.ErrUserAlreadyLoggedIn
	}

	return authService.Login(ctx, email, password)
}

func (r *queryResolver) LoggedUser(ctx context.Context) (*model.User, error) {
	user := users.UserFromContext(ctx)
	if user == nil {
		return nil, errors.ErrUserNotLoggedIn
	}

	return users.FilterUser(user), nil
}
