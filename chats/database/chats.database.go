package chatsDatabase

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Chat struct {
	mgm.DefaultModel   `bson:",inline"`
	ID                 primitive.ObjectID   `bson:"_id"`
	Users              []primitive.ObjectID `bson:"users"`
	Messages           []*Message           `bson:"messages"`
	MessagesCount      uint32               `bson:"messagesCount"`
	LastViewedMessages []LastViewedMessage  `bson:"lastViewedMessages"`
	DeletedBy          []primitive.ObjectID `bson:"deletedBy"`
}

type MessageType struct {
	Text      string
	Reply     string
	File      string
	Files     string
	Image     string
	Images    string
	Sticker   string
	GIF       string
	Withdrawn string
}

var (
	messageType = MessageType{
		Text:      "plaintext",
		Reply:     "reply",
		GIF:       "gif",
		Image:     "image",
		Images:    "images",
		Sticker:   "sticker",
		File:      "file",
		Files:     "files",
		Withdrawn: "withdrawn",
	}
)

type Message struct {
	ID          primitive.ObjectID   `bson:"_id"`
	User        primitive.ObjectID   `bson:"user"`
	Type        string               `bson:"type"`
	Content     string               `bson:"content"`
	SendTime    primitive.DateTime   `bson:"sendTime"`
	Reactions   []*Reaction          `bson:"reactions"`
	IsForwarded bool                 `bson:"isForwarded"`
	DeletedBy   []primitive.ObjectID `bson:"deletedBy"`
}

type Reaction struct {
	Emoji string             `bson:"emoji"`
	User  primitive.ObjectID `bson:"user"`
}

type LastViewedMessage struct {
	Message primitive.ObjectID `bson:"message"`
	User    primitive.ObjectID `bson:"user"`
	Time    primitive.DateTime `bson:"time"`
}

type DB struct{}

// GetDB Returns new DB instance
func GetDB() *DB {
	return &DB{}
}

// UseCollection Returns users *mgm.Collection
func (usersDB *DB) UseCollection() *mgm.Collection {
	return mgm.Coll(&Chat{})
}
