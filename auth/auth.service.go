package auth

import (
	"context"
	"golang.org/x/crypto/bcrypt"
	"messagewith-server/env"
	errors "messagewith-server/error-constants"
	"messagewith-server/graph/model"
	"messagewith-server/sessions"
	"messagewith-server/users"
	"messagewith-server/utils"
)

type service struct{}

func (service *service) Login(ctx context.Context, email string, password string) (*model.User, error) {
	ginCtx := utils.GinContextFromContext(ctx)

	MessagewithJwtSecret := []byte(env.JwtSecret)

	user, err := users.Service.GetPlainUser(ctx, nil, &email, nil)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.ErrUserBadPassword
	}

	session := sessions.Service.CreateSession(ctx, user)
	hash, err := utils.Encrypt(MessagewithJwtSecret, session.Token)
	if err != nil {
		panic(err)
	}

	ginCtx.SetCookie("SessionToken", hash, 60*60*24*7, "/", env.Domain, true, true)

	return users.FilterUser(user), nil
}

func (service *service) Logout(ctx context.Context) (*bool, error) {
	ginCtx := utils.GinContextFromContext(ctx)
	sessionToken, err := sessions.GetSessionTokenFromCookie(ginCtx)
	if err != nil {
		return nil, err
	}

	ok := sessions.Service.ClearSession(ctx, *sessionToken)
	if ok == false {
		return nil, errors.ErrUserNotLoggedIn
	}

	ginCtx.SetCookie("SessionToken", "", 60*60*24*7, "/", env.Domain, true, true)

	return &ok, nil
}
