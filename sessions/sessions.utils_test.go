package sessions

import (
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"messagewith-server/env"
	database "messagewith-server/sessions/database"
	"messagewith-server/users"
	usersDatabase "messagewith-server/users/database"
	"messagewith-server/utils"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetUserFromSession(t *testing.T) {
	id := primitive.NewObjectID()
	usersMockDB := []*usersDatabase.User{
		{ID: primitive.NewObjectID(), FirstName: "Alice", LastName: "Collins"},
		{ID: id, FirstName: "Alice", LastName: "Collins"},
		{ID: primitive.NewObjectID(), FirstName: "Alice", LastName: "Collins"},
	}
	usersMockRepo := new(users.MockRepository)
	usersFindOneRes, usersFindOneErr, usersFindOneHandler := users.GetFindOneRunHandler(&usersMockDB)
	usersMockRepo.On("FindOne", mock.Anything).Run(usersFindOneHandler).Return(&*usersFindOneRes, &*usersFindOneErr)
	users.Service = users.GetService(usersMockRepo)

	user, err := GetUserFromSession(nil, &database.Session{User: id})
	assert.Nil(t, err)
	assert.Equal(t, user, usersMockDB[1])
}

func TestGetSessionTokenFromCookie(t *testing.T) {
	env.InitEnvConstants()
	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)

	firstSession := database.Session{ID: primitive.NewObjectID(), Token: uuid.NewV4().String()}
	encryptedToken, _ := utils.Encrypt([]byte(env.JwtSecret), firstSession.Token)

	r.GET("/", func(ctx *gin.Context) {
		sessionToken, err := GetSessionTokenFromCookie(ctx)
		assert.Nil(t, err)
		assert.Equal(t, firstSession.Token, *sessionToken)
	})

	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{
		Name:     "SessionToken",
		Value:    encryptedToken,
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Date(2999, time.April, 10, 10, 10, 10, 10, time.UTC),
		Path:     "/",
	})

	r.ServeHTTP(w, req)
}

func TestGetLocationFromIP(t *testing.T) {
	env.InitEnvConstants()
	location, err := GetLocationFromIP("168.68.176.64")
	assert.Nil(t, err)
	assert.Equal(t, location.Latitude, 37.751)
	assert.Equal(t, location.Longitude, -97.822)
	assert.Equal(t, location.Country.IsoCode, "US")
	assert.Equal(t, location.Country.IsInEuropeanUnion, false)
	assert.Equal(t, location.AccuracyRadius, uint16(1000))
	assert.Equal(t, location.TimeZone, "America/Chicago")
}
