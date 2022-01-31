package users

import (
	"context"
	"github.com/naamancurtis/mongo-go-struct-to-bson/mapper"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	database "messagewith-server/users/database"
	"reflect"
)

type ResetPasswordMockRepository struct {
	mock.Mock
}

func (m *ResetPasswordMockRepository) FindOne(_ context.Context, filter interface{}) (*database.ResetPassword, error) {
	args := m.Called(filter)
	result := args.Get(0)
	err := args.Get(1).(*error)

	return *result.(**database.ResetPassword), *err
}

func GetResetPasswordFindOneRunHandler(mockDB *[]*database.ResetPassword) (**database.ResetPassword, *error, func(args mock.Arguments)) {
	var res *database.ResetPassword
	var err error

	return &res, &err, func(args mock.Arguments) {
		res = nil
		err = nil
		filters := args.Get(0)

		resetPasswordType := reflect.TypeOf((*database.ResetPassword)(nil)).String()
		filtersType := reflect.TypeOf(filters).String()
		if filtersType == resetPasswordType {
			for _, item := range *mockDB {
				document := filters.(*database.ResetPassword)
				if document == item {
					res = item
				}
			}
		} else {
			filtersMap := filters.(primitive.M)

			for _, item := range *mockDB {
				bsonItem := mapper.ConvertStructToBSONMap(item, nil)
				isEqual := false

				for key, filter := range filtersMap {
					filterKind := reflect.ValueOf(filter).Kind()

					if filterKind == reflect.Ptr {
						if bsonItem[key] == *filter.(*string) {
							isEqual = true
						}
					} else {
						if bsonItem[key] == filter {
							isEqual = true
						}
					}
				}

				if isEqual {
					res = item
				}
			}
		}

		if res == nil {
			err = mongo.ErrNoDocuments
		}
	}
}

func (m *ResetPasswordMockRepository) DeleteOne(_ context.Context, filter interface{}) (*mongo.DeleteResult, error) {
	args := m.Called(filter)
	result := args.Get(0).(**mongo.DeleteResult)
	err := args.Get(1).(*error)

	return *result, *err
}

func GetResetPasswordDeleteOneRunHandler(db *[]*database.ResetPassword) (**mongo.DeleteResult, *error, func(args mock.Arguments)) {
	var (
		res *mongo.DeleteResult
		err error
	)

	return &res, &err, func(args mock.Arguments) {
		res = nil
		err = nil
		filters := args.Get(0)

		var foundIndex = -1

		resetPasswordType := reflect.TypeOf((*database.ResetPassword)(nil)).String()
		filtersType := reflect.TypeOf(filters).String()
		if filtersType == resetPasswordType {
			for index, item := range *db {
				document := filters.(*database.ResetPassword)
				if document == item {
					foundIndex = index
				}
			}
		} else {
			filtersMap := filters.(primitive.M)

			for index, item := range *db {
				bsonItem := mapper.ConvertStructToBSONMap(item, nil)
				isEqual := false

				for key, filter := range filtersMap {
					filterKind := reflect.ValueOf(filter).Kind()

					if filterKind == reflect.Ptr {
						if bsonItem[key] == *filter.(*string) {
							isEqual = true
						}
					} else {
						if bsonItem[key] == filter {
							isEqual = true
						}
					}
				}

				if isEqual {
					foundIndex = index
				}
			}
		}

		if foundIndex == -1 {
			err = mongo.ErrNoDocuments
			return
		}

		(*db)[foundIndex] = (*db)[len(*db)-1]
		(*db)[len(*db)-1] = nil
		*db = (*db)[:len(*db)-1]

		res = &mongo.DeleteResult{}
	}
}

func (m *ResetPasswordMockRepository) Create(document *database.ResetPassword) error {
	args := m.Called(document)
	err := args.Get(0)

	if err == nil {
		return nil
	}

	return err.(error)
}

func GetResetPasswordCreateRunHandler(db *[]*database.ResetPassword) func(args mock.Arguments) {
	return func(args mock.Arguments) {
		document := args.Get(0).(*database.ResetPassword)
		*db = append(*db, document)
	}
}
