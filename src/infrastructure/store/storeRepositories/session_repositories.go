package storeRepositories

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/wataru-dev/bot-api/src/infrastructure/store"
	"github.com/wataru-dev/bot-api/src/infrastructure/store/model"
)

type SessionRepository struct {
	Client *store.FireStoreClient
}

func NewSessionRepository(client *store.FireStoreClient) *SessionRepository{
	return &SessionRepository{
		Client: client,
	}
}

func (sr *SessionRepository) Add(userId, role, content string) error {
	ctx := context.Background()
	_, _, err := sr.Client.Client.Collection("sessions").Doc(userId).Collection("messages").Add(ctx, model.Session{
		Role: role,
		Content: content,
		Timestamp: time.Now().Unix(),
		ExpireAt:  time.Now().Add(24 * time.Hour),
	})
	if err != nil {
		return err
	}
	return nil
}

func (sr *SessionRepository) GetRecentMessages(userId string, limit int) ([]model.Session, error) {
	ctx := context.Background()

	iter := sr.Client.Client.
		Collection("sessions").
		Doc(userId).
		Collection("messages").
		OrderBy("timestamp", firestore.Desc). // 新しい順
		Limit(limit).
		Documents(ctx)

	var messages []model.Session
	for {
		doc, err := iter.Next()
		if err != nil {
			break
		}
		var m model.Session
		if err := doc.DataTo(&m); err != nil {
			return nil, err
		}

		messages = append([]model.Session{m}, messages...)
	}

	return messages, nil
}