package users

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	errors "messagewith-server/error-constants"
	"messagewith-server/graph/model"
	database "messagewith-server/users/database"
	"testing"
)

func TestUserFromContext(t *testing.T) {
	user := &database.User{}
	ctx := context.WithValue(context.Background(), "LoggedUser", user)

	userFromCtx := UserFromContext(ctx)
	assert.Equal(t, userFromCtx, user)

	secondUserFromCtx := UserFromContext(context.Background())
	assert.Nil(t, secondUserFromCtx)
}

func TestCreateNickname(t *testing.T) {
	mockDB := []*database.User{
		{FirstName: "Alice", LastName: "Santiago", Nickname: "alice"},
		{FirstName: "Harrison", LastName: "Hill", Nickname: "harrison_hill"},
		{FirstName: "Harrison", LastName: "Hill", Nickname: "harrison_hill_2"},
		{FirstName: "Otis", LastName: "Dean", Nickname: "otis_dean"},
	}
	mockRepo := new(MockRepository)

	findOneRes, findOneErr, findOneHandler := GetFindOneRunHandler(&mockDB)
	mockRepo.On("FindOne", mock.Anything).Run(findOneHandler).Return(&*findOneRes, &*findOneErr)

	nickname := "alice"
	newNickname, err := createNickname(nil, mockRepo, &model.UserInput{FirstName: "Alice", LastName: "Santiago", Nickname: &nickname})
	assert.Nil(t, newNickname)
	assert.ErrorIs(t, err, errors.ErrUserNicknameAlreadyUsed)

	nickname = "alice2"
	newNickname, err = createNickname(nil, mockRepo, &model.UserInput{FirstName: "Alice", LastName: "Santiago", Nickname: &nickname})
	assert.Equal(t, *newNickname, nickname)
	assert.Nil(t, err)

	newNickname, err = createNickname(nil, mockRepo, &model.UserInput{FirstName: "Harrison", LastName: "Hill"})
	assert.Equal(t, *newNickname, "harrison_hill_3")
	assert.Nil(t, err)

	newNickname, err = createNickname(nil, mockRepo, &model.UserInput{FirstName: "Otis", LastName: "Dean"})
	assert.Equal(t, *newNickname, "otis_dean_2")
	assert.Nil(t, err)

	newNickname, err = createNickname(nil, mockRepo, &model.UserInput{FirstName: "Janusz", LastName: "Kowalski"})
	assert.Equal(t, *newNickname, "janusz_kowalski")
	assert.Nil(t, err)
}
