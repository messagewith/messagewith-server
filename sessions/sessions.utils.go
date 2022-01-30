package sessions

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/oschwald/geoip2-golang"
	"log"
	"messagewith-server/env"
	errors "messagewith-server/error-constants"
	sessionDatabase "messagewith-server/sessions/database"
	"messagewith-server/users"
	database "messagewith-server/users/database"
	"messagewith-server/utils"
	"net"
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

	encryptedSessionToken, err := ctx.Cookie("SessionToken")
	if err != nil || encryptedSessionToken == "" {
		return nil, errors.ErrUserNotLoggedIn
	}

	sessionToken, err := utils.Decrypt(MessagewithJwtSecret, encryptedSessionToken)
	if err != nil {
		return nil, fmt.Errorf("invalid session id: %v", err)
	}

	return &sessionToken, nil
}

func GetLocationFromIP(ip string) (*sessionDatabase.Location, error) {
	geolocationDB, err := geoip2.Open(fmt.Sprintf("%v/geolite/GeoLite2-City.mmdb", env.RootDir))
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
