package sessions

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	ua "github.com/mileusna/useragent"
	"github.com/oschwald/geoip2-golang"
	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"messagewith-server/env"
	errors "messagewith-server/errors"
	"messagewith-server/users"
	"messagewith-server/utils"
	"net"
	"os"
	"time"
)

func GetUserFromSession(ctx context.Context, session *Session) (*users.User, error) {
	userId := session.User.Hex()
	usersService := users.GetService()

	user, err := usersService.GetPlainUser(ctx, &userId, nil, nil)
	if err != nil {
		log.Panicf("failed to get user: %v", err)
	}

	return user, nil
}

func GetSessionTokenFromCookie(ctx *gin.Context) (*string, error) {
	MessagewithJwtSecret := []byte(os.Getenv(env.JwtSecret))

	encryptedSessionToken, err := ctx.Cookie("session_token")
	if err != nil || encryptedSessionToken == "" {
		return nil, errors.ErrUserNotLoggedIn
	}

	sessionToken, err := utils.Decrypt(MessagewithJwtSecret, encryptedSessionToken)
	if err != nil {
		return nil, fmt.Errorf("invalid session id: %v", err)
	}

	return &sessionToken, nil
}

func GetSessionFromCookie(ctx *gin.Context) (*Session, error) {
	sessionToken, err := GetSessionTokenFromCookie(ctx)
	if err != nil {
		return nil, err
	}

	return GetSession(ctx, *sessionToken)
}

func GetLocationFromIP(ip string) (*Location, error) {
	geolocationDB, err := geoip2.Open("geolite/GeoLite2-City.mmdb")
	parsedIp := net.ParseIP(ip)

	if err != nil {
		return nil, err
	}

	defer geolocationDB.Close()

	location := Location{}
	city, err := geolocationDB.City(parsedIp)
	if err != nil {
		return nil, err
	}

	location.Country.IsoCode = city.Country.IsoCode
	location.Country.IsInEuropeanUnion = city.Country.IsInEuropeanUnion
	location.Latitude = city.Location.Latitude
	location.Longitude = city.Location.Longitude
	location.AccuracyRadius = city.Location.AccuracyRadius
	location.TimeZone = city.Location.TimeZone

	return &location, err
}

func CreateSession(ctx context.Context, user *users.User) *Session {
	db := GetDB().UseCollection()
	gc, err := utils.GinContextFromContext(ctx)
	if err != nil {
		panic(err)
	}

	clientIp := gc.ClientIP()
	if clientIp == "::1" {
		clientIp = os.Getenv(env.MockupIpAddress)
	}

	location, err := GetLocationFromIP(clientIp)
	if err != nil {
		panic(err)
	}

	userAgent := gc.GetHeader("User-Agent")
	parsedUserAgent := ua.Parse(userAgent)

	session := Session{
		ID:           primitive.NewObjectID(),
		Token:        uuid.NewV4().String(),
		User:         user.ID,
		OS:           parsedUserAgent.OS,
		Expires:      primitive.NewDateTimeFromTime(time.Now().AddDate(0, 0, 7)),
		LastTimeUsed: primitive.NewDateTimeFromTime(time.Now()),
		Location:     *location,
	}

	err = db.Create(&session)
	if err != nil {
		panic(err)
	}

	return &session
}

func GetSession(ctx context.Context, token string) (*Session, error) {
	db := GetDB().UseCollection()

	result := &Session{}
	err := db.FindOne(ctx, bson.M{"token": token}).Decode(result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func ClearSession(ctx context.Context, token string) bool {
	db := GetDB().UseCollection()

	_, err := db.DeleteOne(ctx, bson.M{"token": token})
	if err != nil {
		return false
	}

	return true
}
