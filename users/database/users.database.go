package usersDatabase

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DB struct{}

type User struct {
	mgm.DefaultModel `bson:",inline"`
	ID               primitive.ObjectID `bson:"_id"`
	FirstName        string             `bson:"firstName"`
	MiddleName       *string            `bson:"middleName,omitempty"`
	LastName         string             `bson:"lastName"`
	FullName         string             `bson:"fullName"`
	Nickname         string             `bson:"nickname"`
	Email            string             `bson:"email"`
	Password         string             `bson:"password"`
}

// GetDB Returns new DB instance
func GetDB() *DB {
	return &DB{}
}

// UseCollection Returns users *mgm.Collection
func (usersDB *DB) UseCollection() *mgm.Collection {
	return mgm.Coll(&User{})
}
