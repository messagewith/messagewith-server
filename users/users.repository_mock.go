package users

import (
	"context"
	"github.com/naamancurtis/mongo-go-struct-to-bson/mapper"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
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

func (m *MockRepository) UpdateByID(_ context.Context, id interface{}, update interface{}) (*mongo.UpdateResult, error) {
	args := m.Called(id, update)
	res := args.Get(0).(**mongo.UpdateResult)
	err := args.Get(1)

	if err == nil {
		return *res, nil
	}

	return *res, *err.(*error)
}

func GetUpdateByIDRunHandler(db *[]*database.User) (**mongo.UpdateResult, *error, func(args mock.Arguments)) {
	var (
		res *mongo.UpdateResult
		err error
	)

	return &res, &err, func(args mock.Arguments) {
		res = nil
		err = nil
		id := args.Get(0).(primitive.ObjectID)
		update := args.Get(1).(primitive.M)

		var updateObjIndex = -1

		for index, item := range *db {
			if item.ID == id {
				updateObjIndex = index
			}
		}

		if updateObjIndex == -1 {
			err = mongo.ErrNoDocuments
			return
		}

		updatedObjBson := mapper.ConvertStructToBSONMap((*db)[updateObjIndex], nil)

		for key, value := range update {
			if updatedObjBson[key] != nil {
				updatedObjBson[key] = value
			}
		}

		var updatedObj *database.User
		bsonBytes, _ := bson.Marshal(updatedObjBson)
		err := bson.Unmarshal(bsonBytes, &updatedObj)
		if err != nil {
			panic(err)
		}
		(*db)[updateObjIndex] = updatedObj

		res = &mongo.UpdateResult{}
	}
}
