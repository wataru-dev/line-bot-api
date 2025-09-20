package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"google.golang.org/genai"

	"github.com/wataru-dev/bot-api/src/config"
	"github.com/wataru-dev/bot-api/src/controller"
	"github.com/wataru-dev/bot-api/src/domain/entities/line"
	"github.com/wataru-dev/bot-api/src/infrastructure/store/model"
)

type BotUseCase struct {
	Repository ISessionRepository
}

type ISessionRepository interface {
	Add(userId, role, content string) error
	GetRecentMessages(userId string, limit int) (*[]model.Session, error)
}


func NewBotUseCase(repository ISessionRepository) controller.IBotUseCase {
	return &BotUseCase{
		Repository: repository,
	}
}

func(buc *BotUseCase) ReplyText(events *line.LineWebhook) error {

	env := config.SetEnvironment()

	for _, e := range events.Events {

		// ユーザーからのメッセージをFireStoreに書き込み
		if err := buc.Repository.Add(e.Source.UserID, "user", e.Message.Text); err != nil {
			return err
		}

		ctx := context.Background()

		// GenAI のクライアントを生成 
		geminiClient, err := genai.NewClient(ctx, &genai.ClientConfig{
        APIKey:  env.GeminiKey,
        Backend: genai.BackendGeminiAPI,
    })

		if err != nil {
				log.Fatal(err)
		}
		
		systemPrompt := "あなたは親しみやすいアシスタントです。ユーザーの入力に対して猫語で答えてください。"

		// トーク履歴を取得
		history, err := buc.Repository.GetRecentMessages(e.Source.UserID, 20)

		if err != nil {
			return err
		}

		prompt := BuildPrompt(systemPrompt, e.Message.Text, history)
		
		// Geminiへのリクエスト
		result, _ := geminiClient.Models.GenerateContent(
        ctx,
        "gemini-2.5-flash",
        genai.Text(prompt),
        nil,
		)

		// リプライ用のペイロードを作成
		payload := line.ReplyMessage{
			ReplyToken: e.ReplyToken,
			Messages: []line.Message{
				{Type: "text", Text: result.Text()},
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

		// トーク履歴をFireStoreに書き込み
		if err := buc.Repository.Add(e.Source.UserID, "assistant", result.Text()); err != nil {
			return err
		}

	}

	return nil

}

func BuildPrompt(systemPrompt, userPrompt string, history *[]model.Session) string {

		var sb strings.Builder

		// システムプロンプトを生成
		if systemPrompt != "" {
			sb.WriteString("### システム指示\n")
			sb.WriteString(systemPrompt + "\n\n")
		}

		// セッション履歴の生成
		if len(*history) > 0 {
			sb.WriteString("### 会話履歴\n")
			for _, m := range *history {
				t := time.Unix(m.Timestamp, 0).Format("2006-01-02 15:04:05")
				switch m.Role {
				case "user":
					sb.WriteString(fmt.Sprintf("[ユーザー @%s]: %s\n", t, m.Content))
				case "assistant":
					sb.WriteString(fmt.Sprintf("[AI @%s]: %s\n", t, m.Content))
				}
			}
			sb.WriteString("\n")
	}

	// ユーザープロンプト
	sb.WriteString("### 新しい入力\n")
	sb.WriteString(userPrompt + "\n")

	return sb.String()
}
