package users

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	database "messagewith-server/users/database"
)

type R interface {
	FindOne(ctx context.Context, filter interface{}) (*database.User, error)
	Find(ctx context.Context, filter interface{}) ([]*database.User, error)
	Create(document *database.User) error
	UpdateByID(ctx context.Context, id interface{}, filter interface{}) (*mongo.UpdateResult, error)
}

type Repository struct{}

func (r *Repository) FindOne(ctx context.Context, filter interface{}) (*database.User, error) {
	user := &database.User{}
	err := collection.FindOne(ctx, filter).Decode(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *Repository) Find(ctx context.Context, filter interface{}) ([]*database.User, error) {
	users := make([]*database.User, 0)
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *Repository) Create(document *database.User) error {
	err := collection.Create(document)

	return err
}

func (r *Repository) UpdateByID(ctx context.Context, id interface{}, update interface{}) (*mongo.UpdateResult, error) {
	return collection.UpdateByID(ctx, id, bson.M{"$set": update})
}
