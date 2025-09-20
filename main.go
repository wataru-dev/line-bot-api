package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/wataru-dev/bot-api/src/controller"
	"github.com/wataru-dev/bot-api/src/infrastructure/web"
)

func main() {

	if _, exists := os.LookupEnv("GOOGLE_CLOUD_PROJECT"); !exists {
		_ = godotenv.Load(".env") // エラーは無視
	}

	engine := web.SetupEngine()

	//	initialize controller
	botController := controller.NewBotController()

	//	routing
	engine.POST("/webhook", botController.Webhook)

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	engine.Run(":" + port)
}
