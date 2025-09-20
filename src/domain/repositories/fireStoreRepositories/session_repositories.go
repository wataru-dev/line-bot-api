package fireStoreRepositories

import (
	"github.com/wataru-dev/bot-api/src/infrastructure/store/model"
	"github.com/wataru-dev/bot-api/src/infrastructure/store/storeRepositories"
	"github.com/wataru-dev/bot-api/src/usecase"
)

type SessionRepository struct {
	Repository *storeRepositories.SessionRepositories
}

func NewSessionRepository(repository *storeRepositories.SessionRepositories) usecase.ISessionRepository {
	return &SessionRepository{
		Repository: repository,
	}
}

func(sr *SessionRepository) Add(userId, role, content string) error {
	if err := sr.Repository.Add(userId, role, content); err != nil {
		return err
	}
	return nil
}

func(sr *SessionRepository) GetRecentMessages(userId string, limit int) (*[]model.Session, error) {
	sessions, err := sr.Repository.GetRecentMessages(userId, 10)

	if err != nil {
		return nil, err
	}
	return &sessions, nil
}