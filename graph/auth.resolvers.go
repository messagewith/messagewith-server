package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"messagewith-server/auth"
	errorConstants "messagewith-server/error-constants"
	"messagewith-server/graph/model"
	"messagewith-server/users"
)

func (r *mutationResolver) Logout(ctx context.Context) (*bool, error) {
	if user := users.GetUserFromContext(ctx); user == nil {
		return nil, errorConstants.ErrUserNotLoggedIn
	}
	return auth.Service.Logout(ctx)
}

func (r *mutationResolver) Login(ctx context.Context, email string, password string) (*model.User, error) {
	if user := users.GetUserFromContext(ctx); user != nil {
		return nil, errorConstants.ErrUserAlreadyLoggedIn
	}
	return auth.Service.Login(ctx, email, password)
}

func (r *queryResolver) LoggedUser(ctx context.Context) (*model.User, error) {
	user := users.GetUserFromContext(ctx)
	if user == nil {
		return nil, errorConstants.ErrUserNotLoggedIn
	}

	return users.FilterUser(user), nil
}
