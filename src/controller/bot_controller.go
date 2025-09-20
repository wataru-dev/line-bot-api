package controller

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wataru-dev/bot-api/src/config"
)

type LineWebhook struct {
	Events []struct {
		Type       string `json:"type"`
		ReplyToken string `json:"replyToken"`
		Source     struct {
			UserID string `json:"userId"`
		} `json:"source"`
		Message struct {
			ID   string `json:"id"`
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"message"`
	} `json:"events"`
}

type ReplyMessage struct {
	ReplyToken string    `json:"replyToken"`
	Messages   []Message `json:"messages"`
}

type Message struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type BotController struct {
}

type IBotUseCase interface {
	ResponseToLine() error
}

func NewBotController() *BotController {
	return &BotController{}
}

func (bc *BotController) Webhook(ctx *gin.Context) {

	//環境変数をロード
	env := config.SetEnvironment()

	var req LineWebhook

	//　バインド
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	//　スライスから必要情報を抽出 → 処理実行
	for _, e := range req.Events {

		// リプライ用のペイロードを作成
		payload := ReplyMessage{
			ReplyToken: e.ReplyToken,
			Messages: []Message{
				{Type: "text", Text: e.Message.Text},
			},
		}

		body, err := json.Marshal(payload)

		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}

		// リクエスト生成
		call, err := http.NewRequest("POST", env.ReplyUri, bytes.NewBuffer(body))
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}
		call.Header.Set("Content-Type", "application/json")
		call.Header.Set("Authorization", "Bearer "+env.LineToken)

		client := &http.Client{}

		// リクエスト実行
		resp, err := client.Do(call)
		if err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
		}
		log.Println(resp)
		defer resp.Body.Close()
	}
}
