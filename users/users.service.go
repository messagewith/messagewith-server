package users

import (
	"context"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	errors "messagewith-server/error-constants"
	"messagewith-server/graph/model"
	"messagewith-server/mails"
	database "messagewith-server/users/database"
	"messagewith-server/utils"
)

type service struct{}

var (
	repository   R
	repositoryRP RP
)

func GetService(rep R, resetPasswordRep RP) *service {
	repository = rep
	repositoryRP = resetPasswordRep
	return &service{}
}

func (service *service) CreateUser(ctx context.Context, userInput *model.UserInput) (*model.User, error) {
	err := validateUserInput(userInput)
	if err != nil {
		return nil, err
	}

	_, err = repository.FindOne(ctx, bson.M{"email": userInput.Email})
	if err == nil {
		return nil, errors.ErrUserEmailAlreadyUsed
	}

	nickname, err := createNickname(ctx, repository, userInput)
	if err != nil {
		return nil, err
	}

	middleName := ""
	if userInput.MiddleName != nil {
		middleName = *userInput.MiddleName + " "
	}

	user := database.User{
		ID:         primitive.NewObjectID(),
		FirstName:  userInput.FirstName,
		MiddleName: userInput.MiddleName,
		LastName:   userInput.LastName,
		FullName:   fmt.Sprintf("%v %v%v", userInput.FirstName, middleName, userInput.LastName),
		Email:      userInput.Email,
		Password:   utils.HashPassword(userInput.Password),
		Nickname:   *nickname,
	}

	err = repository.Create(&user)
	if err != nil {
		panic(err)
	}

	return FilterUser(&user), nil
}

func (service *service) GetUsers(ctx context.Context, filter *model.UserFilter) ([]*model.User, error) {
	filterObj := bson.M{}

	if filter != nil {
		if filter.FirstName != nil {
			filterObj["firstName"] = *filter.FirstName
		}

		if filter.LastName != nil {
			filterObj["lastName"] = *filter.LastName
		}

		if filter.Email != nil {
			filterObj["email"] = *filter.Email
		}

		if filter.FullName != nil {
			filterObj["fullName"] = *filter.FullName
		}
	}

	allUsers, err := repository.Find(ctx, filterObj)
	if err != nil {
		return nil, err
	}

	return FilterAllUsers(allUsers), nil
}

func (service *service) GetUser(ctx context.Context, id *string, email *string, nickname *string) (*model.User, error) {
	user, err := service.GetPlainUser(ctx, id, email, nickname)
	if err != nil {
		return nil, err
	}

	return FilterUser(user), nil
}

func (service *service) GetPlainUser(ctx context.Context, id *string, email *string, nickname *string) (*database.User, error) {
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

	user, err := repository.FindOne(ctx, filterObj)
	if err != nil {
		return nil, possibleErr
	}

	return user, nil
}

func (service *service) GenerateChangePasswordToken(ctx context.Context, email string) (*string, error) {
	user, err := service.GetPlainUser(ctx, nil, &email, nil)
	if err != nil {
		return nil, errors.ErrNoUserWithSpecifiedEmail
	}

	_, err = repositoryRP.FindOne(ctx, bson.M{"user": user.ID})
	if err == nil {
		_, err := repositoryRP.DeleteOne(ctx, bson.M{"user": user.ID})
		if err != nil {
			panic(err)
		}
	}

	token := uuid.NewV4().String()
	resetPasswordDocument := database.ResetPassword{
		ID:    primitive.NewObjectID(),
		Token: token,
		User:  user.ID,
	}

	err = repositoryRP.Create(&resetPasswordDocument)
	if err != nil {
		panic(err)
	}

	ok := mails.Service.SendResetPasswordToken(email, token)
	if ok == false {
		panic("Failed to send reset password e-mail")
	}

	returnMessage := "Check you e-mail inbox"

	return &returnMessage, nil
}

func (service *service) ChangePassword(ctx context.Context, email string, token string, newPassword string) (*model.User, error) {
	res, err := repositoryRP.FindOne(ctx, bson.M{"token": token})
	if err != nil {
		return nil, errors.ErrChangePasswordTokenNotFound
	}

	userId := res.User.Hex()
	user, err := service.GetPlainUser(ctx, &userId, nil, nil)
	if err != nil {
		panic(err)
	}

	if user.Email != email {
		return nil, errors.ErrNoUserWithSpecifiedEmail
	}

	if passwordLength := len(newPassword); passwordLength < 8 {
		return nil, errors.ErrUserPasswordTooShort
	} else if passwordLength > 128 {
		return nil, errors.ErrUserPasswordTooLong
	}

	if !utils.IsPasswordValid(newPassword) {
		return nil, errors.ErrUserBadPassword
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(newPassword))
	if err == nil {
		return nil, errors.ErrChangePasswordSameNewPassword
	}

	user.Password = utils.HashPassword(newPassword)
	_, err = repository.UpdateByID(ctx, user.ID, bson.M{"password": user.Password})
	if err != nil {
		panic(err)
	}

	_, err = repositoryRP.DeleteOne(ctx, res)
	if err != nil {
		panic(err)
	}

	return FilterUser(user), nil
}

func FilterUser(user *database.User) *model.User {
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

func FilterAllUsers(users []*database.User) []*model.User {
	newUsers := make([]*model.User, 0)

	for _, v := range users {
		newUsers = append(newUsers, FilterUser(v))
	}

	return newUsers
}
