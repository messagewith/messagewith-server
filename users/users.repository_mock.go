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

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) FindOne(_ context.Context, filter interface{}) (*database.User, error) {
	args := m.Called(filter)
	result := args.Get(0)
	err := args.Get(1).(*error)

	return *result.(**database.User), *err
}

func GetFindOneRunHandler(mockDB *[]*database.User) (**database.User, *error, func(args mock.Arguments)) {
	var res *database.User
	var err error

	return &res, &err, func(args mock.Arguments) {
		res = nil
		err = nil
		filters := args.Get(0).(primitive.M)

		for _, item := range *mockDB {
			bsonItem := mapper.ConvertStructToBSONMap(item, nil)
			isEqual := false

			for key, filter := range filters {
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

		if res == nil {
			err = mongo.ErrNoDocuments
		}
	}
}

func (m *MockRepository) Find(_ context.Context, filter interface{}) ([]*database.User, error) {
	args := m.Called(filter)
	result := args.Get(0)
	err := args.Get(1).(*error)
	return *result.(*[]*database.User), *err
}

func GetFindRunHandler(mockDB *[]*database.User) (*[]*database.User, *error, func(args mock.Arguments)) {
	res := make([]*database.User, 0)
	var err error

	return &res, &err, func(args mock.Arguments) {
		res = make([]*database.User, 0)
		err = nil

		filters := args.Get(0).(primitive.M)

		for _, item := range *mockDB {
			bsonItem := mapper.ConvertStructToBSONMap(item, nil)
			isEqual := false

			for key, filter := range filters {
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
				res = append(res, item)
			}
		}

		if len(res) == 0 {
			err = mongo.ErrNoDocuments
		}
	}
}

func (m *MockRepository) Create(_ *database.User) error {
	args := m.Called()
	return args.Error(0)
}
