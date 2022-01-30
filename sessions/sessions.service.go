package sessions

import (
	"context"
	"github.com/gin-gonic/gin"
	ua "github.com/mileusna/useragent"
	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"messagewith-server/env"
	"messagewith-server/sessions/database"
	database "messagewith-server/users/database"
	"messagewith-server/utils"
	"time"
)

type service struct{}

var (
	repository R
)

func getService(rep R) *service {
	repository = rep
	return &service{}
}

func (service *service) CreateSession(ctx context.Context, user *database.User) *sessionDatabase.Session {
	gc := utils.GinContextFromContext(ctx)

	clientIp := gc.ClientIP()
	if clientIp == "::1" {
		clientIp = env.MockupIpAddress
	}

	location, err := GetLocationFromIP(clientIp)
	if err != nil {
		panic(err)
	}

	userAgent := gc.GetHeader("User-Agent")
	parsedUserAgent := ua.Parse(userAgent)

	session := sessionDatabase.Session{
		ID:           primitive.NewObjectID(),
		Token:        uuid.NewV4().String(),
		User:         user.ID,
		OS:           parsedUserAgent.OS,
		Expires:      primitive.NewDateTimeFromTime(time.Now().AddDate(0, 0, 7)),
		LastTimeUsed: primitive.NewDateTimeFromTime(time.Now()),
		Location:     *location,
	}

	err = repository.Create(&session)
	if err != nil {
		panic(err)
	}

	return &session
}

func (service *service) GetSession(ctx context.Context, token string) (*sessionDatabase.Session, error) {
	result, err := repository.FindOne(ctx, bson.M{"token": token})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (service *service) GetSessionFromCookie(ctx *gin.Context) (*sessionDatabase.Session, error) {
	sessionToken, err := GetSessionTokenFromCookie(ctx)
	if err != nil {
		return nil, err
	}

	return service.GetSession(ctx, *sessionToken)
}

func (service *service) ClearSession(ctx context.Context, token string) bool {
	_, err := repository.DeleteOne(ctx, bson.M{"token": token})
	if err != nil {
		return false
	}

	return true
}
