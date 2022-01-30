package users

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	errorConstants "messagewith-server/error-constants"
	"messagewith-server/graph/model"
	database "messagewith-server/users/database"
	"testing"
)

func TestService_CreateUser(t *testing.T) {
	mockDB := []*database.User{{
		ID:        primitive.NewObjectID(),
		Email:     "johny@email.com",
		FirstName: "Johny",
		LastName:  "Miller",
		Nickname:  "johny_miller",
	}}
	testObj := new(MockRepository)

	findOneResult, findOneErr, findOneHandler := GetFindOneRunHandler(&mockDB)
	testObj.On("FindOne", mock.Anything).Run(findOneHandler).Return(&*findOneResult, &*findOneErr)
	testObj.On("Create").Return(nil)
	service := getService(testObj)

	_, err := service.CreateUser(nil, &model.UserInput{})
	assert.ErrorIs(t, err, errorConstants.ErrUserInputNotContainsAllProps)

	// First name too short
	_, err = service.CreateUser(nil, &model.UserInput{FirstName: "Y", LastName: "as", Email: "johny@email.com", Password: "asgyib@!S3"})
	assert.ErrorIs(t, err, errorConstants.ErrUserFirstNameTooShort)

	// First name too long
	_, err = service.CreateUser(nil, &model.UserInput{FirstName: "Yaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaas", LastName: "as", Email: "email@email.com", Password: "asgyib@!S3"})
	assert.ErrorIs(t, err, errorConstants.ErrUserFirstNameTooLong)

	// First name good
	_, err = service.CreateUser(nil, &model.UserInput{FirstName: "as", LastName: "as", Email: "johny@email.com", Password: "asgyib@!S33"})
	assert.ErrorIs(t, err, errorConstants.ErrUserEmailAlreadyUsed)

	// Last name too short
	_, err = service.CreateUser(nil, &model.UserInput{FirstName: "Yolo", LastName: "a", Email: "johny@email.com", Password: "asgyib@!S3"})
	assert.ErrorIs(t, err, errorConstants.ErrUserLastNameTooShort)

	// Last name too long
	_, err = service.CreateUser(nil, &model.UserInput{FirstName: "Yolo", LastName: "Yaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", Email: "johny@email.com", Password: "asgyib@!S3"})
	assert.ErrorIs(t, err, errorConstants.ErrUserLastNameTooLong)

	// Last name good
	_, err = service.CreateUser(nil, &model.UserInput{FirstName: "Yolo", LastName: "as", Email: "johny@email.com", Password: "asgyib@!S3"})
	assert.ErrorIs(t, err, errorConstants.ErrUserEmailAlreadyUsed)

	// E-mail too short
	_, err = service.CreateUser(nil, &model.UserInput{FirstName: "Yolo", LastName: "Yolo", Email: "asd", Password: "asgyib@!S3"})
	assert.ErrorIs(t, err, errorConstants.ErrUserEmailTooShort)

	// E-mail too long
	_, err = service.CreateUser(nil, &model.UserInput{FirstName: "Yolo", LastName: "Yolo", Email: "YBo1DYsUOEsLpKo7crnBUqJVX6Gg9VHpEADk9Qy3T3chqFZfkBwj4DnTcTJcXPZiDKbHLLloDaJSxpivlNqc9K5x2RhLlLafg5b4aG8Ex7rz6XBcdceN07LcrJdv1pTbymbPodDBf9J0XeuvdZtVCDmM1VI53ovRNeGxBk2puhyzVpTTc3qzp8CcHXBzZ36PFL1afedkhsiWOFHbfaixWDD6ODB5TcOBgVqWR9LjvKKMG210fgaZXnlvAb3mX1j", Password: "asgyib@!S3"})
	assert.ErrorIs(t, err, errorConstants.ErrUserEmailTooLong)

	// E-mail invalid
	_, err = service.CreateUser(nil, &model.UserInput{FirstName: "Yolo", LastName: "Yolo", Email: "asdddddd", Password: "asgyib@!S3"})
	assert.ErrorIs(t, err, errorConstants.ErrUserEmailWrong)

	// E-mail good
	_, err = service.CreateUser(nil, &model.UserInput{FirstName: "Yolo", LastName: "Yolo", Email: "johny@email.com", Password: "asgyib@!S3"})
	assert.ErrorIs(t, err, errorConstants.ErrUserEmailAlreadyUsed)

	// Nickname too short
	nickname := "a"
	_, err = service.CreateUser(nil, &model.UserInput{FirstName: "Yolo", LastName: "Yolo", Email: "johny@email.com", Password: "asgyib@!S3", Nickname: &nickname})
	assert.ErrorIs(t, err, errorConstants.ErrUserNicknameTooShort)

	// Nickname too long
	nickname = "aaaaaaaaaaaaaaaaaaaaaaaaaaa"
	_, err = service.CreateUser(nil, &model.UserInput{FirstName: "Yolo", LastName: "Yolo", Email: "johny@email.com", Password: "asgyib@!S3", Nickname: &nickname})
	assert.ErrorIs(t, err, errorConstants.ErrUserNicknameTooLong)

	// Nickname good
	nickname = "aaaaaaaaaaaaaaaaaaaaaaaaaa"
	_, err = service.CreateUser(nil, &model.UserInput{FirstName: "Yolo", LastName: "Yolo", Email: "johny@email.com", Password: "asgyib@!S3", Nickname: &nickname})
	assert.ErrorIs(t, err, errorConstants.ErrUserEmailAlreadyUsed)

	// Middle name too short
	middleName := "a"
	_, err = service.CreateUser(nil, &model.UserInput{FirstName: "Yolo", LastName: "Yolo", Email: "johny@email.com", Password: "asgyib@!S3", MiddleName: &middleName})
	assert.ErrorIs(t, err, errorConstants.ErrUserMiddleNameTooShort)

	// Middle name too long
	middleName = "Yaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaas"
	_, err = service.CreateUser(nil, &model.UserInput{FirstName: "Yolo", LastName: "Yolo", Email: "johny@email.com", Password: "asgyib@!S3", MiddleName: &middleName})
	assert.ErrorIs(t, err, errorConstants.ErrUserMiddleNameTooLong)

	// Middle name good
	middleName = "to"
	_, err = service.CreateUser(nil, &model.UserInput{FirstName: "Yolo", LastName: "Yolo", Email: "johny@email.com", Password: "asgyib@!S3", MiddleName: &middleName})
	assert.ErrorIs(t, err, errorConstants.ErrUserEmailAlreadyUsed)

	// Password too short
	_, err = service.CreateUser(nil, &model.UserInput{FirstName: "Yolo", LastName: "Yolo", Email: "johny@email.com", Password: "assss"})
	assert.ErrorIs(t, err, errorConstants.ErrUserPasswordTooShort)

	// Password too long
	_, err = service.CreateUser(nil, &model.UserInput{FirstName: "Yolo", LastName: "Yolo", Email: "johny@email.com", Password: "epTjLKTmHagT5ulTIwtAViIubLDZ48XZQBE9xBMf6rQicVxqRzg59qanbnMAZPloV27Nx1NXlQ3Qf3UL1umdOlzjoNOas4wBB4MJZRBnYchi3kBmhyUNiS6ci9eEvAMb9"})
	assert.ErrorIs(t, err, errorConstants.ErrUserPasswordTooLong)

	// Password invalid
	_, err = service.CreateUser(nil, &model.UserInput{FirstName: "Yolo", LastName: "Yolo", Email: "johny@email.com", Password: "ssssssss"})
	assert.ErrorIs(t, err, errorConstants.ErrUserBadPassword)
	_, err = service.CreateUser(nil, &model.UserInput{FirstName: "Yolo", LastName: "Yolo", Email: "johny@email.com", Password: "ssssssssS"})
	assert.ErrorIs(t, err, errorConstants.ErrUserBadPassword)
	_, err = service.CreateUser(nil, &model.UserInput{FirstName: "Yolo", LastName: "Yolo", Email: "johny@email.com", Password: "ssssssssS1"})
	assert.ErrorIs(t, err, errorConstants.ErrUserBadPassword)

	// Password good
	_, err = service.CreateUser(nil, &model.UserInput{FirstName: "Yolo", LastName: "Yolo", Email: "johny@email.com", Password: "ssssssssS1!"})
	assert.ErrorIs(t, err, errorConstants.ErrUserEmailAlreadyUsed)

	// If user with received nickname exists, func should return error
	nickname = "johny_miller"
	_, err = service.CreateUser(nil, &model.UserInput{FirstName: "Yolo", LastName: "Yolo", Email: "another@email.com", Password: "ssssssssS1!", Nickname: &nickname})
	assert.ErrorIs(t, err, errorConstants.ErrUserNicknameAlreadyUsed)

	createdUser, _ := service.CreateUser(nil, &model.UserInput{FirstName: "Johny", LastName: "Rambo", Email: "another@email.com", Password: "ssssssssS1!"})
	assert.Equal(t, createdUser.Nickname, "johny_rambo")
	assert.Equal(t, createdUser.FullName, "Johny Rambo")
	assert.Equal(t, createdUser.FirstName, "Johny")
	assert.Equal(t, createdUser.LastName, "Rambo")
	assert.Equal(t, createdUser.Email, "another@email.com")

	middleName = "Kowalski"
	createdUser, _ = service.CreateUser(nil, &model.UserInput{FirstName: "Johny", MiddleName: &middleName, LastName: "Miller", Email: "another@email.com", Password: "ssssssssS1!"})
	assert.Equal(t, createdUser.Nickname, "johny_miller_2")
	assert.Equal(t, createdUser.FullName, "Johny Kowalski Miller")
	assert.Equal(t, createdUser.FirstName, "Johny")
	assert.Equal(t, createdUser.LastName, "Miller")
	assert.Equal(t, createdUser.Email, "another@email.com")

	testObj.AssertExpectations(t)
}

func TestService_GetUsers(t *testing.T) {
	var mockDB = make([]*database.User, 0)
	mockDB = append(mockDB, &database.User{FirstName: "Jakub", LastName: "Kowalski", Email: "jakub@kowalski.com", FullName: "Jakub Kowalski"})
	mockDB = append(mockDB, &database.User{FirstName: "Jakub", LastName: "Rambo", Email: "jakub@rambo.com", FullName: "Jakub Rambo"})
	mockDB = append(mockDB, &database.User{FirstName: "Andrzej", LastName: "Kowalski", Email: "andrzej@kowalski2.com", FullName: "Andrzej Kowalski"})

	testObj := new(MockRepository)
	findRes, findErr, findHandler := GetFindRunHandler(&mockDB)
	testObj.On("Find", mock.AnythingOfType("primitive.M")).Run(findHandler).Return(&*findRes, &*findErr)

	service := getService(testObj)
	firstName := "Jakub"
	users, _ := service.GetUsers(nil, &model.UserFilter{FirstName: &firstName})
	assert.Equal(t, users, FilterAllUsers([]*database.User{mockDB[0], mockDB[1]}))

	lastName := "Kowalski"
	users, _ = service.GetUsers(nil, &model.UserFilter{LastName: &lastName})
	assert.Equal(t, users, FilterAllUsers([]*database.User{mockDB[0], mockDB[2]}))

	email := "jakub@kowalski.com"
	users, _ = service.GetUsers(nil, &model.UserFilter{Email: &email})
	assert.Equal(t, users, FilterAllUsers([]*database.User{mockDB[0]}))

	fullName := "Andrzej Kowalski"
	users, _ = service.GetUsers(nil, &model.UserFilter{FullName: &fullName})
	assert.Equal(t, users, FilterAllUsers([]*database.User{mockDB[2]}))

	testObj.AssertExpectations(t)
}

func TestService_GetUser(t *testing.T) {
	id := primitive.NewObjectID()
	hexId := id.Hex()

	mockDB := []*database.User{
		{ID: id, FirstName: "Johny", Email: "johny@rambo.com", Nickname: "johny"},
		{ID: primitive.NewObjectID(), FirstName: "Mark", Email: "mark@rambo.com", Nickname: "mark"},
		{ID: primitive.NewObjectID(), FirstName: "Noname", Email: "noname@rambo.com", Nickname: "noname"},
	}

	testObj := new(MockRepository)
	findOneResult, findOneErr, findOneHandler := GetFindOneRunHandler(&mockDB)
	testObj.On("FindOne", mock.Anything).Run(findOneHandler).Return(&*findOneResult, &*findOneErr)

	service := getService(testObj)
	user, _ := service.GetUser(nil, &hexId, nil, nil)
	assert.Equal(t, user, FilterUser(mockDB[0]))

	email := "mark@rambo.com"
	user, _ = service.GetUser(nil, nil, &email, nil)
	assert.Equal(t, user, FilterUser(mockDB[1]))

	nickname := "noname"
	user, _ = service.GetUser(nil, nil, nil, &nickname)
	assert.Equal(t, user, FilterUser(mockDB[2]))
}

func TestFilterUser(t *testing.T) {
	id := primitive.NewObjectID()
	testUser := &database.User{
		ID:         id,
		FirstName:  "Johny",
		LastName:   "Rambo",
		Email:      "johny@rambo.com",
		Password:   "VerySecretPassword@2!",
		FullName:   "Johny Rambo",
		MiddleName: nil,
		Nickname:   "johny_rambo",
	}

	filteredUser := FilterUser(testUser)
	assert.Equal(t, filteredUser, &model.User{
		ID:         id.Hex(),
		FirstName:  testUser.FirstName,
		LastName:   testUser.LastName,
		FullName:   testUser.FullName,
		Email:      testUser.Email,
		MiddleName: testUser.MiddleName,
		Nickname:   "johny_rambo",
	})
}

func TestFilterAllUser(t *testing.T) {
	id := primitive.NewObjectID()
	id2 := primitive.NewObjectID()
	testAllUsers := []*database.User{
		{
			ID:         id,
			FirstName:  "Johny",
			LastName:   "Rambo",
			Email:      "johny@rambo.com",
			Password:   "VerySecretPassword@2!",
			FullName:   "Johny Rambo",
			MiddleName: nil,
			Nickname:   "johny_rambo",
		},
		{
			ID:         id2,
			FirstName:  "Johny",
			LastName:   "Rambo",
			Email:      "johny@rambo.com",
			Password:   "VerySecretPassword@2!",
			FullName:   "Johny Rambo",
			MiddleName: nil,
			Nickname:   "johny_rambo",
		},
	}
	filteredUsers := FilterAllUsers(testAllUsers)
	assert.Equal(t, filteredUsers, []*model.User{
		{
			ID:         id.Hex(),
			FirstName:  testAllUsers[0].FirstName,
			LastName:   testAllUsers[0].LastName,
			FullName:   testAllUsers[0].FullName,
			Email:      testAllUsers[0].Email,
			MiddleName: testAllUsers[0].MiddleName,
			Nickname:   testAllUsers[0].Nickname,
		},
		{
			ID:         id2.Hex(),
			FirstName:  testAllUsers[1].FirstName,
			LastName:   testAllUsers[1].LastName,
			FullName:   testAllUsers[1].FullName,
			Email:      testAllUsers[1].Email,
			MiddleName: testAllUsers[1].MiddleName,
			Nickname:   testAllUsers[1].Nickname,
		},
	})
}
