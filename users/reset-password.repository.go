package users

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	database "messagewith-server/users/database"
)

type RP interface {
	FindOne(ctx context.Context, filter interface{}) (*database.ResetPassword, error)
	DeleteOne(ctx context.Context, filter interface{}) (*mongo.DeleteResult, error)
	Create(document *database.ResetPassword) error
}

type ResetPasswordRepository struct{}

func (r *ResetPasswordRepository) Create(document *database.ResetPassword) error {
	err := resetPasswordCollection.Create(document)

	return err
}

func (r *ResetPasswordRepository) FindOne(ctx context.Context, filter interface{}) (*database.ResetPassword, error) {
	res := &database.ResetPassword{}
	err := resetPasswordCollection.FindOne(ctx, filter).Decode(res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *ResetPasswordRepository) DeleteOne(ctx context.Context, filter interface{}) (*mongo.DeleteResult, error) {
	res, err := resetPasswordCollection.DeleteOne(ctx, filter)

	return res, err
}
