package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"messagewith-server/chats"
	errorConstants "messagewith-server/error-constants"
	"messagewith-server/graph/generated"
	"messagewith-server/graph/model"
	"messagewith-server/users"
)

func (r *mutationResolver) CreateChat(ctx context.Context, userID string) (*model.Chat, error) {
	user := users.GetUserFromContext(ctx)
	if user == nil {
		return nil, errorConstants.ErrUserNotLoggedIn
	}
	return chats.Service.CreateChat(ctx, user, userID)
}

func (r *mutationResolver) DeleteChat(ctx context.Context, id string) (*bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) SendMessage(ctx context.Context, chatID string, messageType model.MessageType, content string) (*model.Message, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeleteMessage(ctx context.Context, id string, chatID string) (*bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UndoMessage(ctx context.Context, id string, chatID string) (*bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) SetMessageReaction(ctx context.Context, chatID string, messageID string, emoji string) (*model.Reaction, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeleteMessageReaction(ctx context.Context, chatID string, messageID string) (*bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) ChatByID(ctx context.Context, id string) (*model.Chat, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Chats(ctx context.Context, filter *model.ChatFilter) ([]*model.Chat, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *subscriptionResolver) ChatByID(ctx context.Context, id string) (<-chan *model.Chat, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *subscriptionResolver) Chats(ctx context.Context, filter *model.ChatFilter) (<-chan []*model.Chat, error) {
	//id := uuid.NewV4().String()
	//cookies := utils.GinContextFromContext(ctx).Request.Cookies()
	//msgs := make(chan []*model.Chat, 1)
	//
	//// Start a goroutine to allow for cleaning up subscriptions that are disconnected.
	//// This go routine will only get past Done() when a client terminates the subscription. This allows us
	//// to only then remove the reference from the list of ChatObservers since it is no longer needed.
	//go func() {
	//	<-ctx.Done()
	//	r.mu.Lock()
	//	delete(r.ChatObservers, id)
	//	r.mu.Unlock()
	//}()
	//r.mu.Lock()
	//// Keep a reference of the channel so that we can push changes into it when new messages are posted.
	//r.ChatObservers[id] = msgs
	//r.mu.Unlock()
	//// This is optional, and this allows newly subscribed clients to get a list of all the messages that have been
	//// posted so far. Upon subscribing the client will be pushed the messages once, further changes are handled
	//// in the PostMessage mutation.
	//r.ChatObservers[id] <- r.ChatMessages
	//return msgs, nil
	panic("uniplemented")
}

// Subscription returns generated.SubscriptionResolver implementation.
func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }

type subscriptionResolver struct{ *Resolver }
