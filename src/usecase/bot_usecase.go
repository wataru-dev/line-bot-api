package usecase

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/wataru-dev/bot-api/src/config"
	"github.com/wataru-dev/bot-api/src/controller"
	"github.com/wataru-dev/bot-api/src/domain/entities"
)

type BotUseCase struct {
}


func NewBotUseCase() controller.IBotUseCase {
	return &BotUseCase{

	}
}

func(buc *BotUseCase) ReplyText(events *entities.LineWebhook) error {

	env := config.SetEnvironment()

	for _, e := range events.Events {
		// リプライ用のペイロードを作成
		payload := entities.ReplyMessage{
			ReplyToken: e.ReplyToken,
			Messages: []entities.Message{
				{Type: "text", Text: e.Message.Text},
			},
		}

		body, err := json.Marshal(payload)

		if err != nil {
			return err
		}

		// リクエスト生成
		call, err := http.NewRequest("POST", env.ReplyUri, bytes.NewBuffer(body))
		if err != nil {
			return err
		}
		call.Header.Set("Content-Type", "application/json")
		call.Header.Set("Authorization", "Bearer "+env.LineToken)

		client := &http.Client{}

		// リクエスト実行
		resp, err := client.Do(call)
		if err != nil {
			return err
		}
		log.Println(resp)

		defer resp.Body.Close()

	}

	return nil

}
