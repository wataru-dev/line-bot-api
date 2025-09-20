package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/wataru-dev/bot-api/src/controller"
	"github.com/wataru-dev/bot-api/src/domain/repositories/fireStoreRepositories"
	"github.com/wataru-dev/bot-api/src/infrastructure/store"
	"github.com/wataru-dev/bot-api/src/infrastructure/store/storeRepositories"
	"github.com/wataru-dev/bot-api/src/infrastructure/web"
	"github.com/wataru-dev/bot-api/src/usecase"
)

func main() {

	if _, exists := os.LookupEnv("GOOGLE_CLOUD_PROJECT"); !exists {
		_ = godotenv.Load(".env") // エラーは無視
	}

	engine := web.SetupEngine()

	// initialize client
	fireStoreClient, _ := store.NewFireStoreClient()
	defer fireStoreClient.Close()

	// initialize store repository
	infrastructureSessionRepository := storeRepositories.NewSessionRepository(fireStoreClient)

	// initialize repository
	sessionRepository := fireStoreRepositories.NewSessionRepository(infrastructureSessionRepository)

	// initialize usecase
	botUseCase := usecase.NewBotUseCase(sessionRepository)

	//	initialize controller
	botController := controller.NewBotController(botUseCase)

	//	routing
	engine.POST("/webhook", botController.Webhook)

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	engine.Run(":" + port)
}
