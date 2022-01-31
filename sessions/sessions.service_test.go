package sessions

import (
	"context"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math"
	"messagewith-server/env"
	errors "messagewith-server/error-constants"
	database "messagewith-server/sessions/database"
	usersDatabase "messagewith-server/users/database"
	"messagewith-server/utils"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestService_CreateSession(t *testing.T) {
	env.InitEnvConstants()
	w := httptest.NewRecorder()
	mockRepository := new(MockRepository)
	mockRepository.On("Create").Return(nil)
	service := getService(mockRepository)

	testingUserId := primitive.NewObjectID()
	testingUser := usersDatabase.User{ID: testingUserId, FirstName: "Hello", LastName: "World"}

	_, e := gin.CreateTestContext(w)
	e.Use(func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), "GinContextKey", c)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	})

	e.GET("/", func(c *gin.Context) {
		c.ClientIP()
		session := service.CreateSession(c.Request.Context(), &testingUser)
		assert.Equal(t, session.User, testingUserId)
		assert.IsType(t, primitive.ObjectID{}, session.ID)
		_, err := uuid.FromString(session.Token)
		assert.Nil(t, err)
		assert.IsType(t, "", session.OS)
		assert.IsType(t, primitive.NewDateTimeFromTime(time.Now().AddDate(0, 0, 7)), session.Expires)
		assert.IsType(t, primitive.NewDateTimeFromTime(time.Now()), session.LastTimeUsed)
		assert.IsType(t, database.Location{}, session.Location)
		duration := session.Expires.Time().Sub(time.Now())
		hours := int(math.Round(duration.Hours()))
		assert.Equal(t, 168, hours)
		duration = session.LastTimeUsed.Time().Sub(time.Now())
		hours = int(math.Round(duration.Hours()))
		assert.Equal(t, 0, hours)
	})

	req := httptest.NewRequest("GET", "/", nil)
	e.ServeHTTP(w, req)
}

func TestService_GetSession(t *testing.T) {
	token := uuid.NewV4().String()
	mockDB := []*database.Session{
		{ID: primitive.NewObjectID(), Token: uuid.NewV4().String()},
		{ID: primitive.NewObjectID(), Token: token},
		{ID: primitive.NewObjectID(), Token: uuid.NewV4().String()},
	}
	mockRepository := new(MockRepository)
	findOneRes, findOneErr, findOneHandler := GetFindOneRunHandler(&mockDB)
	mockRepository.On("FindOne", mock.Anything).Run(findOneHandler).Return(&*findOneRes, &*findOneErr)

	service := getService(mockRepository)
	session, err := service.GetSession(nil, token)
	assert.Nil(t, err)
	assert.Equal(t, mockDB[1], session)

	session, err = service.GetSession(nil, uuid.NewV4().String())
	assert.Nil(t, session)
	assert.NotNil(t, err)
}

func TestService_GetSessionFromCookie(t *testing.T) {
	env.InitEnvConstants()
	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)

	token := uuid.NewV4().String()
	encryptedToken, _ := utils.Encrypt([]byte(env.JwtSecret), token)

	mockDB := []*database.Session{
		{ID: primitive.NewObjectID(), Token: uuid.NewV4().String()},
		{ID: primitive.NewObjectID(), Token: token},
		{ID: primitive.NewObjectID(), Token: uuid.NewV4().String()},
	}

	mockRepository := new(MockRepository)
	findOneRes, findOneErr, findOneHandler := GetFindOneRunHandler(&mockDB)
	mockRepository.On("FindOne", mock.Anything).Run(findOneHandler).Return(&*findOneRes, &*findOneErr)
	service := getService(mockRepository)

	r.GET("/", func(c *gin.Context) {
		session, err := service.GetSessionFromCookie(c)
		assert.Nil(t, err)
		assert.Equal(t, mockDB[1], session)
	})

	r.GET("/second-req", func(c *gin.Context) {
		session, err := service.GetSessionFromCookie(c)
		assert.Nil(t, session)
		assert.ErrorIs(t, errors.ErrUserNotLoggedIn, err)
	})

	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{
		Name:     "SessionToken",
		Value:    strings.ReplaceAll(encryptedToken, "+", "%2B"),
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Date(2999, time.April, 10, 10, 10, 10, 10, time.UTC),
		Path:     "/",
	})

	secondReq := httptest.NewRequest("GET", "/second-req", nil)
	r.ServeHTTP(w, req)
	r.ServeHTTP(w, secondReq)
}
