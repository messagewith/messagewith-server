package users

import (
	"context"
	"fmt"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	errors "messagewith-server/errors"
	"messagewith-server/graph/model"
	"strings"
)

func UserFromContext(ctx context.Context) *User {
	user, _ := ctx.Value("LoggedUser").(*User)
	return user
}

func createNickname(ctx context.Context, db *mgm.Collection, userInput *model.UserInput) (*string, error) {
	foundUser := &User{}
	if userInput.Nickname != nil {
		if err := db.FindOne(ctx, bson.M{"nickname": userInput.Nickname}).Decode(foundUser); err == nil {
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
			tempNickname = fmt.Sprintf("%v_%v", firstNameAndLastName, i)
		}

		if err := db.FindOne(ctx, bson.M{"nickname": tempNickname}).Decode(foundUser); err != nil {
			newNickname = &tempNickname
		}

		i++
	}

	return newNickname, nil
}
