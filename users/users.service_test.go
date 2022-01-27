package users

import (
	"context"
	"errors"
	"github.com/naamancurtis/mongo-go-struct-to-bson/mapper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	errorConstants "messagewith-server/error-constants"
	"messagewith-server/graph/model"
	database "messagewith-server/users/database"
	"testing"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) FindOne(_ context.Context, _ interface{}) (*database.User, error) {
	args := m.Called()
	result := args.Get(0)

	if result != nil {
		return result.(*database.User), args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *MockService) Find(_ context.Context, filter interface{}) ([]*database.User, error) {
	args := m.Called(filter)
	result := args.Get(0)
	return *result.(*[]*database.User), args.Error(1)
}

func (m *MockService) Create(_ *database.User) error {
	args := m.Called()
	return args.Error(0)
}

func TestCreateUser(t *testing.T) {
	testObj := new(MockService)
	findOneResult := &database.User{
		ID:        primitive.NewObjectID(),
		Email:     "johny@email.com",
		FirstName: "Johny",
		LastName:  "Miller",
		Nickname:  "johny_miller",
	}
	testObj.On("FindOne").Return(findOneResult, nil).Times(6)
	testObj.On("Create").Return(nil)
	service := getService(testObj)

	_, err := service.CreateUser(nil, &model.UserInput{})
	assert.ErrorIs(t, err, errorConstants.ErrUserInputNotContainsAllProps)

	// First name too short
	_, err = service.CreateUser(nil, &model.UserInput{FirstName: "Y", LastName: "as", Email: findOneResult.Email, Password: "asgyib@!S3"})
	assert.ErrorIs(t, err, errorConstants.ErrUserFirstNameTooShort)

	// First name too long
	_, err = service.CreateUser(nil, &model.UserInput{FirstName: "Yaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaas", LastName: "as", Email: "email@email.com", Password: "asgyib@!S3"})
	assert.ErrorIs(t, err, errorConstants.ErrUserFirstNameTooLong)

	// First name good
	_, err = service.CreateUser(nil, &model.UserInput{FirstName: "as", LastName: "as", Email: findOneResult.Email, Password: "asgyib@!S33"})
	assert.ErrorIs(t, err, errorConstants.ErrUserEmailAlreadyUsed)

	// Last name too short
	_, err = service.CreateUser(nil, &model.UserInput{FirstName: "Yolo", LastName: "a", Email: findOneResult.Email, Password: "asgyib@!S3"})
	assert.ErrorIs(t, err, errorConstants.ErrUserLastNameTooShort)

	// Last name too long
	_, err = service.CreateUser(nil, &model.UserInput{FirstName: "Yolo", LastName: "Yaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", Email: findOneResult.Email, Password: "asgyib@!S3"})
	assert.ErrorIs(t, err, errorConstants.ErrUserLastNameTooLong)

	// Last name good
	_, err = service.CreateUser(nil, &model.UserInput{FirstName: "Yolo", LastName: "as", Email: findOneResult.Email, Password: "asgyib@!S3"})
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
	_, err = service.CreateUser(nil, &model.UserInput{FirstName: "Yolo", LastName: "Yolo", Email: findOneResult.Email, Password: "asgyib@!S3"})
	assert.ErrorIs(t, err, errorConstants.ErrUserEmailAlreadyUsed)

	// Nickname too short
	nickname := "a"
	_, err = service.CreateUser(nil, &model.UserInput{FirstName: "Yolo", LastName: "Yolo", Email: findOneResult.Email, Password: "asgyib@!S3", Nickname: &nickname})
	assert.ErrorIs(t, err, errorConstants.ErrUserNicknameTooShort)

	// Nickname too long
	nickname = "aaaaaaaaaaaaaaaaaaaaaaaaaaa"
	_, err = service.CreateUser(nil, &model.UserInput{FirstName: "Yolo", LastName: "Yolo", Email: findOneResult.Email, Password: "asgyib@!S3", Nickname: &nickname})
	assert.ErrorIs(t, err, errorConstants.ErrUserNicknameTooLong)

	// Nickname good
	nickname = "aaaaaaaaaaaaaaaaaaaaaaaaaa"
	_, err = service.CreateUser(nil, &model.UserInput{FirstName: "Yolo", LastName: "Yolo", Email: findOneResult.Email, Password: "asgyib@!S3", Nickname: &nickname})
	assert.ErrorIs(t, err, errorConstants.ErrUserEmailAlreadyUsed)

	// Middle name too short
	middleName := "a"
	_, err = service.CreateUser(nil, &model.UserInput{FirstName: "Yolo", LastName: "Yolo", Email: findOneResult.Email, Password: "asgyib@!S3", MiddleName: &middleName})
	assert.ErrorIs(t, err, errorConstants.ErrUserMiddleNameTooShort)

	// Middle name too long
	middleName = "Yaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaas"
	_, err = service.CreateUser(nil, &model.UserInput{FirstName: "Yolo", LastName: "Yolo", Email: findOneResult.Email, Password: "asgyib@!S3", MiddleName: &middleName})
	assert.ErrorIs(t, err, errorConstants.ErrUserMiddleNameTooLong)

	// Middle name good
	middleName = "to"
	_, err = service.CreateUser(nil, &model.UserInput{FirstName: "Yolo", LastName: "Yolo", Email: findOneResult.Email, Password: "asgyib@!S3", MiddleName: &middleName})
	assert.ErrorIs(t, err, errorConstants.ErrUserEmailAlreadyUsed)

	// Password too short
	_, err = service.CreateUser(nil, &model.UserInput{FirstName: "Yolo", LastName: "Yolo", Email: findOneResult.Email, Password: "assss"})
	assert.ErrorIs(t, err, errorConstants.ErrUserPasswordTooShort)

	// Password too long
	_, err = service.CreateUser(nil, &model.UserInput{FirstName: "Yolo", LastName: "Yolo", Email: findOneResult.Email, Password: "epTjLKTmHagT5ulTIwtAViIubLDZ48XZQBE9xBMf6rQicVxqRzg59qanbnMAZPloV27Nx1NXlQ3Qf3UL1umdOlzjoNOas4wBB4MJZRBnYchi3kBmhyUNiS6ci9eEvAMb9"})
	assert.ErrorIs(t, err, errorConstants.ErrUserPasswordTooLong)

	// Password invalid
	_, err = service.CreateUser(nil, &model.UserInput{FirstName: "Yolo", LastName: "Yolo", Email: findOneResult.Email, Password: "ssssssss"})
	assert.ErrorIs(t, err, errorConstants.ErrUserBadPassword)
	_, err = service.CreateUser(nil, &model.UserInput{FirstName: "Yolo", LastName: "Yolo", Email: findOneResult.Email, Password: "ssssssssS"})
	assert.ErrorIs(t, err, errorConstants.ErrUserBadPassword)
	_, err = service.CreateUser(nil, &model.UserInput{FirstName: "Yolo", LastName: "Yolo", Email: findOneResult.Email, Password: "ssssssssS1"})
	assert.ErrorIs(t, err, errorConstants.ErrUserBadPassword)

	// Password good
	_, err = service.CreateUser(nil, &model.UserInput{FirstName: "Yolo", LastName: "Yolo", Email: findOneResult.Email, Password: "ssssssssS1!"})
	assert.ErrorIs(t, err, errorConstants.ErrUserEmailAlreadyUsed)

	testObj.On("FindOne").Return(nil, errors.New("")).Once()
	testObj.On("FindOne").Return(findOneResult, nil).Once()
	// If user with received nickname exists, func should return error
	_, err = service.CreateUser(nil, &model.UserInput{FirstName: "Yolo", LastName: "Yolo", Email: "another@email.com", Password: "ssssssssS1!", Nickname: &findOneResult.Nickname})
	assert.ErrorIs(t, err, errorConstants.ErrUserNicknameAlreadyUsed)

	testObj.On("FindOne").Return(nil, errors.New("")).Twice()
	createdUser, _ := service.CreateUser(nil, &model.UserInput{FirstName: "Johny", LastName: "Rambo", Email: "another@email.com", Password: "ssssssssS1!"})
	assert.Equal(t, createdUser.Nickname, "johny_rambo")
	assert.Equal(t, createdUser.FullName, "Johny Rambo")
	assert.Equal(t, createdUser.FirstName, "Johny")
	assert.Equal(t, createdUser.LastName, "Rambo")
	assert.Equal(t, createdUser.Email, "another@email.com")

	testObj.On("FindOne").Return(nil, errors.New("")).Once()
	testObj.On("FindOne").Return(findOneResult, nil).Once()
	testObj.On("FindOne").Return(nil, errors.New("")).Once()
	middleName = "Kowalski"
	createdUser, _ = service.CreateUser(nil, &model.UserInput{FirstName: "Johny", MiddleName: &middleName, LastName: "Miller", Email: "another@email.com", Password: "ssssssssS1!"})
	assert.Equal(t, createdUser.Nickname, "johny_miller_1")
	assert.Equal(t, createdUser.FullName, "Johny Kowalski Miller")
	assert.Equal(t, createdUser.FirstName, "Johny")
	assert.Equal(t, createdUser.LastName, "Miller")
	assert.Equal(t, createdUser.Email, "another@email.com")

	testObj.AssertExpectations(t)
}

func TestGetUsers(t *testing.T) {
	var mockDB = make([]*database.User, 0)
	mockDB = append(mockDB, &database.User{FirstName: "Jakub", LastName: "Kowalski", Email: "jakub@kowalski.com", FullName: "Jakub Kowalski"})
	mockDB = append(mockDB, &database.User{FirstName: "Jakub", LastName: "Rambo", Email: "jakub@rambo.com", FullName: "Jakub Rambo"})
	mockDB = append(mockDB, &database.User{FirstName: "Andrzej", LastName: "Kowalski", Email: "andrzej@kowalski2.com", FullName: "Andrzej Kowalski"})
	result := make([]*database.User, 0)

	testObj := new(MockService)
	testObj.On("Find", mock.AnythingOfType("primitive.M")).Run(func(args mock.Arguments) {
		result = make([]*database.User, 0)
		filters := args.Get(0).(primitive.M)

		for _, item := range mockDB {
			itemBsonMap := mapper.ConvertStructToBSONMap(item, nil)
			badItem := false

			for key, filter := range filters {
				if itemBsonMap[key] != filter {
					badItem = true
				}
			}

			if !badItem {
				result = append(result, item)
			}
		}
	}).Return(&result, nil)

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
