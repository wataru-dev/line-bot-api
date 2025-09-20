package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wataru-dev/bot-api/src/domain/entities"
)



type BotController struct {
	UseCase IBotUseCase
}

type IBotUseCase interface {
	ReplyText(events *entities.LineWebhook) error
}

func NewBotController(useCase IBotUseCase) *BotController {
	return &BotController{
		UseCase: useCase,
	}
}

func (bc *BotController) Webhook(ctx *gin.Context) {

	var req entities.LineWebhook

	//　バインド
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	//　スライスから必要情報を抽出 → 処理実行
	for _, e := range req.Events {

		switch e.Type {
		case "text":
			err := bc.UseCase.ReplyText(&req)
			if err != nil {
				ctx.JSON(400, gin.H{"error": err.Error()})
			} else {
				ctx.JSON(200, gin.H{"message": "reply success"})
			}
		}
	}

}
