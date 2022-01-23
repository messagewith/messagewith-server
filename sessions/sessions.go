package sessions

import (
	"awesomeProject/users"
	"awesomeProject/utils"
	"context"
	ua "github.com/mileusna/useragent"
	"github.com/oschwald/geoip2-golang"
	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net"
	"os"
	"time"
)

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

func CreateSession(ctx context.Context, user *users.User) (*Session, error) {
	db := GetDB(ctx).UseCollection()
	gc, err := utils.GinContextFromContext(ctx)

	if err != nil {
		return nil, err
	}

	clientIp := gc.ClientIP()

	if clientIp == "::1" {
		clientIp = os.Getenv("MESSAGEWITH_MOCKUP_IP_ADDRESS")
	}

	location, err := GetLocationFromIP(clientIp)

	if err != nil {
		return nil, err
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
		return nil, err
	}

	return &session, nil
}

func GetSession(ctx context.Context, token string) (*Session, error) {
	db := GetDB(ctx).UseCollection()

	result := &Session{}

	err := db.FindOne(ctx, bson.M{"token": token}).Decode(result)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func ClearSession(ctx context.Context, token string) bool {
	db := GetDB(ctx).UseCollection()

	_, err := db.DeleteOne(ctx, bson.M{"token": token})

	if err != nil {
		return false
	}

	return true
}
