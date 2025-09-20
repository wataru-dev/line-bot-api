package store

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/wataru-dev/bot-api/src/config"
)

type FireStoreClient struct {
	Client *firestore.Client
}

// NewFireStoreClient Firestoreクライアントを初期化
func NewFireStoreClient() (*FireStoreClient, error) {
	env := config.SetEnvironment()
	ctx := context.Background()

	// サービスアカウントキーのパスを環境変数などから取得
	// 例: GOOGLE_APPLICATION_CREDENTIALS=/path/to/service-account.json
	client, err := firestore.NewClient(ctx, env.GoogleProjectID)
	if err != nil {
		return nil, fmt.Errorf("failed to create firestore client: %w", err)
	}

	return &FireStoreClient{Client: client}, nil
}

// Close Firestore接続を終了
func (f *FireStoreClient) Close() {
	f.Client.Close()
}
