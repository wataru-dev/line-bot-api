package controller

import (
    "github.com/gin-gonic/gin"
    entities "github.com/wataru-dev/bot-api/src/domain/entities/line"
)
type BotController struct {
    UseCase IBotUseCase
}

type IBotUseCase interface {
    ReplyText(event *entities.LineEvent) error
}

func NewBotController(useCase IBotUseCase) *BotController {
	return &BotController{
		UseCase: useCase,
	}
}

func (bc *BotController) Webhook(ctx *gin.Context) {

    var req entities.LineWebhook

    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(400, gin.H{"error": err.Error()})
        return
    }

    // 1イベントずつ処理、HTTPレスポンスは最後に1回
    var firstErr error
    for i := range req.Events {
        e := &req.Events[i]
        if e.Type != "message" {
            continue
        }
        if err := bc.UseCase.ReplyText(e); err != nil && firstErr == nil {
            firstErr = err
        }
    }

    if firstErr != nil {
        ctx.JSON(400, gin.H{"error": firstErr.Error()})
        return
    }

    ctx.JSON(200, gin.H{"message": "accepted"})
}
