package users

import (
	"context"
	"fmt"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	errors "messagewith-server/errors"
	"messagewith-server/graph/model"
	"messagewith-server/utils"
)

type Service struct {
	db *mgm.Collection
}

func GetService() *Service {
	return &Service{
		db: GetDB().UseCollection(),
	}
}

func (service *Service) CreateUser(ctx context.Context, userInput *model.UserInput) (*model.User, error) {
	foundUser := &User{}
	err := service.db.FindOne(ctx, bson.M{"email": userInput.Email}).Decode(foundUser)
	if err == nil {
		return nil, errors.ErrUserEmailAlreadyUsed
	}

	nickname, err := createNickname(ctx, service.db, userInput)
	if err != nil {
		return nil, err
	}

	middleName := ""
	if userInput.MiddleName != nil {
		middleName = *userInput.MiddleName + " "
	}

	user := User{
		ID:         primitive.NewObjectID(),
		FirstName:  userInput.FirstName,
		MiddleName: userInput.MiddleName,
		LastName:   userInput.LastName,
		FullName:   fmt.Sprintf("%v %v%v", userInput.FirstName, middleName, userInput.LastName),
		Email:      userInput.Email,
		Password:   utils.HashPassword(userInput.Password),
		Nickname:   *nickname,
	}

	err = service.db.Create(&user)
	if err != nil {
		panic(err)
	}

	return FilterUser(&user), nil
}

func (service *Service) GetUsers(ctx context.Context, filter *model.UserFilter) ([]*model.User, error) {
	allUsers := make([]*User, 0)
	filterObj := bson.M{}

	if filter != nil {
		if filter.FirstName != nil {
			filterObj["firstName"] = filter.FirstName
		}

		if filter.LastName != nil {
			filterObj["lastName"] = filter.LastName
		}

		if filter.Email != nil {
			filterObj["email"] = filter.Email
		}
	}

	cursor, err := service.db.Find(ctx, filterObj)
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &allUsers)
	if err != nil {
		return nil, err
	}

	return FilterAllUsers(allUsers), nil
}

func (service *Service) GetUser(ctx context.Context, id *string, email *string, nickname *string) (*model.User, error) {
	user, err := service.GetPlainUser(ctx, id, email, nickname)
	if err != nil {
		return nil, err
	}

	return FilterUser(user), nil
}

func (service *Service) GetPlainUser(ctx context.Context, id *string, email *string, nickname *string) (*User, error) {
	filterObj := bson.M{}
	var possibleErr error

	if id != nil {
		objectId, err := primitive.ObjectIDFromHex(*id)
		if err != nil {
			return nil, errors.ErrInvalidID
		}

		filterObj["_id"] = objectId
		possibleErr = errors.ErrNoUserWithSpecifiedId
	}

	if email != nil {
		filterObj["email"] = email
		possibleErr = errors.ErrNoUserWithSpecifiedEmail
	}

	if nickname != nil {
		filterObj["nickname"] = nickname
		possibleErr = errors.ErrNoUserWithSpecifiedNickname
	}

	res := &User{}
	err := service.db.FindOne(ctx, filterObj).Decode(res)
	if err != nil {
		return nil, possibleErr
	}

	return res, nil
}

func FilterUser(user *User) *model.User {
	return &model.User{
		ID:         user.ID.Hex(),
		FirstName:  user.FirstName,
		MiddleName: user.MiddleName,
		FullName:   user.FullName,
		LastName:   user.LastName,
		Email:      user.Email,
		Nickname:   user.Nickname,
	}
}

func FilterAllUsers(users []*User) []*model.User {
	newUsers := make([]*model.User, 0)

	for _, v := range users {
		newUsers = append(newUsers, FilterUser(v))
	}

	return newUsers
}
