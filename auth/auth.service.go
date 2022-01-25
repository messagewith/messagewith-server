package auth

import (
	"context"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"messagewith-server/env"
	errors "messagewith-server/errors"
	"messagewith-server/graph/model"
	"messagewith-server/sessions"
	"messagewith-server/users"
	"messagewith-server/utils"

	"os"
)

type Service struct {
	usersService *users.Service
	ginCtx       *gin.Context
}

func GetService(ctx context.Context) *Service {
	ginCtx, err := utils.GinContextFromContext(ctx)
	if err != nil {
		panic(err)
	}

	service := &Service{
		usersService: users.GetService(),
		ginCtx:       ginCtx,
	}

	return service
}

func (service *Service) Login(ctx context.Context, email string, password string) (*model.User, error) {
	MessagewithJwtSecret := []byte(os.Getenv(env.JwtSecret))

	user, err := service.usersService.GetPlainUser(ctx, nil, &email, nil)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.ErrUserBadPassword
	}

	session := sessions.CreateSession(ctx, user)
	hash, err := utils.Encrypt(MessagewithJwtSecret, session.Token)
	if err != nil {
		panic(err)
	}

	service.ginCtx.SetCookie("session_token", hash, 60*60*24*7, "/", os.Getenv(env.Domain), true, true)

	return users.FilterUser(user), nil
}

func (service *Service) Logout(ctx context.Context) (*bool, error) {
	sessionToken, err := sessions.GetSessionTokenFromCookie(service.ginCtx)
	if err != nil {
		return nil, err
	}

	ok := sessions.ClearSession(ctx, *sessionToken)
	if ok == false {
		return nil, errors.ErrUserNotLoggedIn
	}

	service.ginCtx.SetCookie("session_token", "", 60*60*24*7, "/", os.Getenv(env.Domain), true, true)

	return &ok, nil
}
