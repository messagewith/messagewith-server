package sessions

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	database "messagewith-server/sessions/database"
)

type R interface {
	FindOne(ctx context.Context, filter interface{}) (*database.Session, error)
	Find(ctx context.Context, filter interface{}) ([]*database.Session, error)
	Create(document *database.Session) error
	DeleteOne(ctx context.Context, filter interface{}) (*mongo.DeleteResult, error)
}

type Repository struct{}

func (r *Repository) FindOne(ctx context.Context, filter interface{}) (*database.Session, error) {
	user := &database.Session{}
	err := collection.FindOne(ctx, filter).Decode(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *Repository) Find(ctx context.Context, filter interface{}) ([]*database.Session, error) {
	users := make([]*database.Session, 0)
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *Repository) Create(document *database.Session) error {
	err := collection.Create(document)

	return err
}

func (r *Repository) DeleteOne(ctx context.Context, filter interface{}) (*mongo.DeleteResult, error) {
	res, err := collection.DeleteOne(ctx, filter)

	return res, err
}
