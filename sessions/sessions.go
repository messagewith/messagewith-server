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
	errors "messagewith-server/error-constants"
	"messagewith-server/sessions/database"
	"messagewith-server/users"
	database "messagewith-server/users/database"
	"messagewith-server/utils"
	"net"
	"time"
)

func GetUserFromSession(ctx context.Context, session *sessionDatabase.Session) (*database.User, error) {
	userId := session.User.Hex()
	user, err := users.Service.GetPlainUser(ctx, &userId, nil, nil)
	if err != nil {
		log.Panicf("failed to get user: %v", err)
	}

	return user, nil
}

func GetSessionTokenFromCookie(ctx *gin.Context) (*string, error) {
	MessagewithJwtSecret := []byte(env.JwtSecret)

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

func GetSessionFromCookie(ctx *gin.Context) (*sessionDatabase.Session, error) {
	sessionToken, err := GetSessionTokenFromCookie(ctx)
	if err != nil {
		return nil, err
	}

	return GetSession(ctx, *sessionToken)
}

func GetLocationFromIP(ip string) (*sessionDatabase.Location, error) {
	geolocationDB, err := geoip2.Open("geolite/GeoLite2-City.mmdb")
	parsedIp := net.ParseIP(ip)

	if err != nil {
		return nil, err
	}

	defer func(geolocationDB *geoip2.Reader) {
		err := geolocationDB.Close()
		if err != nil {
			panic(err)
		}
	}(geolocationDB)

	location := sessionDatabase.Location{}
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

func CreateSession(ctx context.Context, user *database.User) *sessionDatabase.Session {
	db := sessionDatabase.GetDB().UseCollection()
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

	err = db.Create(&session)
	if err != nil {
		panic(err)
	}

	return &session
}

func GetSession(ctx context.Context, token string) (*sessionDatabase.Session, error) {
	db := sessionDatabase.GetDB().UseCollection()

	result := &sessionDatabase.Session{}
	err := db.FindOne(ctx, bson.M{"token": token}).Decode(result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func ClearSession(ctx context.Context, token string) bool {
	db := sessionDatabase.GetDB().UseCollection()

	_, err := db.DeleteOne(ctx, bson.M{"token": token})
	if err != nil {
		return false
	}

	return true
}
