package sessions

import (
	"context"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Location struct {
	Country        Country
	Latitude       float64
	Longitude      float64
	AccuracyRadius uint16
	TimeZone       string
}

type Country struct {
	IsoCode           string
	IsInEuropeanUnion bool
}

type Session struct {
	mgm.DefaultModel `bson:",inline"`
	ID               primitive.ObjectID `bson:"_id"`
	Token            string             `bson:"token"`
	User             primitive.ObjectID `bson:"user"`
	Location         Location           `bson:"location"`
	OS               string             `bson:"os"`
	LastTimeUsed     primitive.DateTime `bson:"lastTimeUsed"`
	Expires          primitive.DateTime `bson:"expires"`
}

type DB struct{}

func GetDB(ctx context.Context) *DB {
	return &DB{}
}

func (db *DB) UseCollection() *mgm.Collection {
	return mgm.Coll(&Session{})
}
