package auth

import (
	"awesomeProject/graph/model"
	"awesomeProject/sessions"
	"awesomeProject/users"
	"awesomeProject/utils"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
	"os"
)

type Service struct{}

func getSessionTokenFromCookie(ctx *gin.Context) (*string, error) {
	MessagewithJwtSecret := []byte(os.Getenv("MESSAGEWITH_JWT_SECRET"))

	encryptedSessionToken, err := ctx.Cookie("session_token")

	if err != nil || encryptedSessionToken == "" {
		return nil, fmt.Errorf("user is not logged in")
	}

	sessionToken, err := utils.Decrypt(MessagewithJwtSecret, encryptedSessionToken)

	if err != nil {
		return nil, fmt.Errorf("invalid session id: %v", err)
	}

	return &sessionToken, nil
}

func getSessionFromCookie(ctx *gin.Context) (*sessions.Session, error) {
	sessionToken, err := getSessionTokenFromCookie(ctx)

	if err != nil {
		return nil, err
	}

	return sessions.GetSession(ctx, *sessionToken)
}

func (service *Service) Login(ctx context.Context, email string, password string) (*model.User, error) {
	usersDB := users.GetDB(ctx).UseCollection()
	ginCtx, err := utils.GinContextFromContext(ctx)
	MessagewithJwtSecret := []byte(os.Getenv("MESSAGEWITH_JWT_SECRET"))

	if err != nil {
		panic(fmt.Errorf("failed to get gin context"))
	}

	_, err = getSessionFromCookie(ginCtx)

	if err == nil {
		return nil, fmt.Errorf("user is already logged")
	}

	user := &users.User{}
	err = usersDB.FindOne(ctx, bson.M{"email": email}).Decode(user)

	if err != nil {
		return nil, fmt.Errorf("can not find user with this e-mail")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		return nil, fmt.Errorf("bad password")
	}

	session, err := sessions.CreateSession(ctx, user)

	if err != nil {
		panic(fmt.Errorf("failed to create session"))
	}

	hash, err := utils.Encrypt(MessagewithJwtSecret, session.Token)

	if err != nil {
		panic(fmt.Errorf("failed to encrypt session token"))
	}

	ginCtx.SetCookie("session_token", hash, 60*60*24*7, "/", os.Getenv("MESSAGEWITH_DOMAIN"), true, true)

	return users.FilterUser(user), nil
}

func (service *Service) Logout(ctx context.Context) (*bool, error) {
	ginCtx, err := utils.GinContextFromContext(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get gin context: %v", err)
	}

	sessionToken, err := getSessionTokenFromCookie(ginCtx)

	if err != nil {
		return nil, err
	}

	ok := sessions.ClearSession(ctx, *sessionToken)

	if ok == false {
		return nil, fmt.Errorf("user is to logged in")
	}

	ginCtx.SetCookie("session_token", "", 60*60*24*7, "/", os.Getenv("MESSAGEWITH_DOMAIN"), true, true)

	return &ok, nil
}

func (service *Service) GetLoggedUser(ctx context.Context) (*model.User, error) {
	ginCtx, err := utils.GinContextFromContext(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get gin context: %v", err)
	}

	session, err := getSessionFromCookie(ginCtx)

	if err != nil {
		return nil, fmt.Errorf("failed to get session: %v", err)
	}

	userId := session.User.Hex()

	usersService := &users.Service{}
	user, err := usersService.GetUser(ctx, &userId, nil)

	if err != nil {
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	return user, nil
}
