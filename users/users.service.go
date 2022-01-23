package users

import (
	"awesomeProject/graph/model"
	"awesomeProject/utils"
	"context"
	"fmt"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
)

type Service struct{}

func createNickname(ctx context.Context, db *mgm.Collection, userInput *model.UserInput) (*string, error) {
	foundUser := &User{}

	if userInput.Nickname != nil {
		if err := db.FindOne(ctx, bson.M{"nickname": userInput.Nickname}).Decode(foundUser); err == nil {
			return nil, fmt.Errorf("user with this nickname already exists")
		}

		return userInput.Nickname, nil
	}

	var (
		newNickname  *string
		tempNickname string
		i            uint = 0
	)
	firstNameAndLastName := fmt.Sprintf("%v_%v", strings.ToLower(userInput.FirstName), strings.ToLower(userInput.LastName))

	for newNickname == nil {
		if i == 0 {
			tempNickname = firstNameAndLastName
		} else {
			tempNickname = fmt.Sprintf("%v_%v", firstNameAndLastName, i)
		}

		if err := db.FindOne(ctx, bson.M{"nickname": tempNickname}).Decode(foundUser); err != nil {
			newNickname = &tempNickname
		}

		i++
	}

	return newNickname, nil
}

func (service *Service) CreateUser(ctx context.Context, userInput *model.UserInput) (*model.User, error) {
	db := GetDB(ctx).UseCollection()

	foundUser := &User{}
	err := db.FindOne(ctx, bson.M{"email": userInput.Email}).Decode(foundUser)

	if err == nil {
		return nil, fmt.Errorf("user with this email already exists")
	}

	nickname, err := createNickname(ctx, db, userInput)

	if err != nil {
		return nil, err
	}

	user := User{
		ID:        primitive.NewObjectID(),
		FirstName: userInput.FirstName,
		LastName:  userInput.LastName,
		Email:     userInput.Email,
		Password:  utils.GeneratePassword(userInput.Password),
		Nickname:  *nickname,
	}

	err = db.Create(&user)

	if err != nil {
		panic(err)
	}

	return FilterUser(&user), nil
}

func (service *Service) GetUsers(ctx context.Context, filter *model.UserFilter) ([]*model.User, error) {
	db := GetDB(ctx).UseCollection()

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

	cursor, err := db.Find(ctx, filterObj)

	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &allUsers)

	if err != nil {
		return nil, err
	}

	return FilterAllUsers(allUsers), nil
}

func (service *Service) GetUser(ctx context.Context, id *string, email *string) (*model.User, error) {
	db := GetDB(ctx).UseCollection()

	filterObj := bson.M{}

	if id != nil {
		objectId, err := primitive.ObjectIDFromHex(*id)

		if err != nil {
			return nil, fmt.Errorf("invalid id")
		}

		filterObj["_id"] = objectId
	}
	if email != nil {
		filterObj["email"] = email
	}

	res := &User{}
	err := db.FindOne(ctx, filterObj).Decode(res)

	if err != nil {
		return nil, err
	}

	return FilterUser(res), nil
}

func FilterUser(user *User) *model.User {
	return &model.User{
		ID:        user.ID.Hex(),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Nickname:  user.Nickname,
	}
}

func FilterAllUsers(users []*User) []*model.User {
	newUsers := make([]*model.User, 0)

	for _, v := range users {
		newUsers = append(newUsers, &model.User{
			ID:        v.ID.String(),
			FirstName: v.FirstName,
			LastName:  v.LastName,
			Email:     v.LastName,
			Nickname:  v.Nickname,
		})
	}

	return newUsers
}
