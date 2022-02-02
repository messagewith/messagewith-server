package chats

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	database "messagewith-server/chats/database"
	errorConstants "messagewith-server/error-constants"
	"messagewith-server/graph/model"
	"messagewith-server/users"
	usersDatabase "messagewith-server/users/database"
)

var (
	repository  R
	cachedUsers = map[string]*model.User{}
)

type service struct{}

func getService(repo R) *service {
	repository = repo
	return &service{}
}

func (service *service) CreateChat(ctx context.Context, user *usersDatabase.User, secondUserID string) (*model.Chat, error) {
	secondUserObjectID, err := primitive.ObjectIDFromHex(secondUserID)
	if err != nil {
		return nil, errorConstants.ErrInvalidID
	}

	_, err = repository.FindOne(ctx, bson.M{"users": bson.M{"$size": 2, "$in": []primitive.ObjectID{user.ID, secondUserObjectID}}})
	if err == nil {
		return nil, errorConstants.ErrChatAlreadyCreated
	}

	secondUser, err := users.Service.GetPlainUser(ctx, &secondUserID, nil, nil)
	if err != nil {
		return nil, errorConstants.ErrInvalidID
	}

	chat := &database.Chat{
		ID:                 primitive.NewObjectID(),
		DeletedBy:          []primitive.ObjectID{},
		LastViewedMessages: []database.LastViewedMessage{},
		Messages:           []*database.Message{},
		MessagesCount:      0,
		Users:              []primitive.ObjectID{user.ID, secondUserObjectID},
	}

	err = repository.Create(ctx, chat)
	if err != nil {
		panic(err)
	}

	return FilterChat(ctx, chat, user.ID, secondUser), nil
}

func FilterChat(ctx context.Context, chat *database.Chat, userObjectID primitive.ObjectID, secondUser *usersDatabase.User) *model.Chat {
	modelChat := &model.Chat{}
	modelChat.MessagesCount = int(chat.MessagesCount)
	modelChat.User = users.FilterUser(secondUser)
	modelChat.ID = chat.ID.Hex()
	modelChat.Messages = FilterAllMessages(ctx, userObjectID, chat.Messages)
	modelChat.LastViewedMessage = FilterLastViewedMessage(ctx, chat.LastViewedMessages, chat.Messages)

	return modelChat
}

func GetChatUser(ctx context.Context, userID primitive.ObjectID) *model.User {
	var filteredUser *model.User
	itemIdHex := userID.Hex()

	if cachedUsers[itemIdHex] != nil {
		filteredUser = cachedUsers[itemIdHex]
	} else {
		user, _ := users.Service.GetUser(ctx, &itemIdHex, nil, nil)
		cachedUsers[itemIdHex] = user
		filteredUser = user
	}

	return filteredUser
}

func GetMessage(allMessages []*database.Message, id primitive.ObjectID) *database.Message {
	for _, item := range allMessages {
		if item.ID == id {
			return item
		}
	}

	return nil
}

func FilterLastViewedMessage(ctx context.Context, lastViewedMessages []database.LastViewedMessage, allMessages []*database.Message) *model.LastViewedMessage {
	if len(lastViewedMessages) == 1 {
		return &model.LastViewedMessage{
			Message: FilterMessage(ctx, nil, GetMessage(allMessages, lastViewedMessages[0].Message)),
			User:    GetChatUser(ctx, lastViewedMessages[0].User),
			Time:    int(lastViewedMessages[0].Time.Time().Unix()),
		}
	}

	return nil
}

func FilterMessage(ctx context.Context, userObjectID *primitive.ObjectID, message *database.Message) *model.Message {
	for _, deletedBy := range message.DeletedBy {
		if deletedBy == *userObjectID {
			return nil
		}
	}

	var filteredUser *model.User
	messageIdHex := message.ID.Hex()

	return &model.Message{
		ID:          messageIdHex,
		User:        filteredUser,
		MessageType: model.MessageType(message.Type),
		Content:     message.Content,
		IsForwarded: message.IsForwarded,
		Reactions:   FilterAllReactions(ctx, message.Reactions),
		SendTime:    int(message.SendTime.Time().Unix()),
	}
}

func FilterAllMessages(ctx context.Context, userObjectID primitive.ObjectID, messages []*database.Message) []*model.Message {
	allModelMessages := make([]*model.Message, 0)

	for _, item := range messages {
		filteredMessage := FilterMessage(ctx, &userObjectID, item)
		if filteredMessage == nil {
			continue
		}

		allModelMessages = append(allModelMessages, filteredMessage)
	}

	return allModelMessages
}

func FilterAllReactions(ctx context.Context, reactions []*database.Reaction) []*model.Reaction {
	allModelReactions := make([]*model.Reaction, 0)

	for _, item := range reactions {
		filteredUser := GetChatUser(ctx, item.User)

		allModelReactions = append(allModelReactions, &model.Reaction{
			User:  filteredUser,
			Emoji: item.Emoji,
		})
	}

	return allModelReactions
}
