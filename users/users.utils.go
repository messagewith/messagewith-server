package users

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	errors "messagewith-server/error-constants"
	"messagewith-server/graph/model"
	database "messagewith-server/users/database"
	"messagewith-server/utils"
	"strings"
)

func UserFromContext(ctx context.Context) *database.User {
	user, _ := ctx.Value("LoggedUser").(*database.User)

	return user
}

func createNickname(ctx context.Context, repository R, userInput *model.UserInput) (*string, error) {
	if userInput.Nickname != nil {
		if _, err := repository.FindOne(ctx, bson.M{"nickname": userInput.Nickname}); err == nil {
			return nil, errors.ErrUserNicknameAlreadyUsed
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
			tempNickname = fmt.Sprintf("%v_%v", firstNameAndLastName, i+1)
		}

		if _, err := repository.FindOne(ctx, bson.M{"nickname": tempNickname}); err != nil {
			newNickname = &tempNickname
		}

		i++
	}

	return newNickname, nil
}

func validateUserInput(input *model.UserInput) error {
	if input == nil || input.FirstName == "" || input.LastName == "" || input.Password == "" {
		return errors.ErrUserInputNotContainsAllProps
	}

	if firstNameLength := len(input.FirstName); firstNameLength < 2 {
		return errors.ErrUserFirstNameTooShort
	} else if firstNameLength > 50 {
		return errors.ErrUserFirstNameTooLong
	}

	if lastNameLength := len(input.LastName); lastNameLength < 2 {
		return errors.ErrUserLastNameTooShort
	} else if lastNameLength > 50 {
		return errors.ErrUserLastNameTooLong
	}

	if input.MiddleName != nil {
		if middleNameLength := len(*input.MiddleName); middleNameLength < 2 {
			return errors.ErrUserMiddleNameTooShort
		} else if middleNameLength > 50 {
			return errors.ErrUserMiddleNameTooLong
		}
	}

	if input.Nickname != nil {
		if nicknameLength := len(*input.Nickname); nicknameLength < 2 {
			return errors.ErrUserNicknameTooShort
		} else if nicknameLength > 26 {
			return errors.ErrUserNicknameTooLong
		}
	}

	if emailLength := len(input.Email); emailLength < 6 {
		return errors.ErrUserEmailTooShort
	} else if emailLength > 254 {
		return errors.ErrUserEmailTooLong
	}

	if !utils.IsEmailValid(input.Email) {
		return errors.ErrUserEmailWrong
	}

	if passwordLength := len(input.Password); passwordLength < 8 {
		return errors.ErrUserPasswordTooShort
	} else if passwordLength > 128 {
		return errors.ErrUserPasswordTooLong
	}

	if !utils.IsPasswordValid(input.Password) {
		return errors.ErrUserBadPassword
	}

	return nil
}
