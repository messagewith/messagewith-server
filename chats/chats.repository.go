package chats

import (
	"context"
	database "messagewith-server/chats/database"
)

type R interface {
	FindOne(ctx context.Context, filter interface{}) (*database.Chat, error)
	Create(ctx context.Context, document *database.Chat) error
}

type Repository struct{}

func (r *Repository) FindOne(ctx context.Context, filter interface{}) (*database.Chat, error) {
	chat := &database.Chat{}
	err := collection.FindOne(ctx, filter).Decode(chat)
	if err != nil {
		return nil, err
	}

	return chat, nil
}

func (r *Repository) Create(ctx context.Context, document *database.Chat) error {
	err := collection.CreateWithCtx(ctx, document)

	return err
}
