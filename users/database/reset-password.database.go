package usersDatabase

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ResetPasswordDB struct{}

type ResetPassword struct {
	mgm.DefaultModel `bson:",inline"`
	ID               primitive.ObjectID `bson:"_id"`
	User             primitive.ObjectID `bson:"user"`
	Token            string             `bson:"token"`
}

func GetResetPasswordDB() *ResetPasswordDB {
	return &ResetPasswordDB{}
}

func (db *ResetPasswordDB) UseCollection() *mgm.Collection {
	return mgm.Coll(&ResetPassword{})
}
