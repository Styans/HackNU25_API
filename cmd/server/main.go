package main

import (
	"ai-assistant/internal/api"
	"ai-assistant/internal/config"
	"ai-assistant/internal/integration/llm"
	"ai-assistant/internal/integration/rerank"
	"ai-assistant/internal/repository"
	"ai-assistant/internal/service"
	"context"
	"log"
	"net/http"
)

func main() {
	ctx := context.Background()

	// 1. Конфиг
	cfg := config.Load()

	// 2. База данных
	dbPool := repository.ConnectDB(ctx, cfg.DatabaseURL)
	defer dbPool.Close()

	// 3. Репозитории (Слой данных)
	userRepo := repository.NewUserRepository(dbPool)
	accRepo := repository.NewAccountRepository(dbPool)
	txRepo := repository.NewTransactionRepository(dbPool)
	finRepo := repository.NewFinancingRepository(dbPool)
	ragRepo := repository.NewRAGRepository(dbPool)
	chatRepo := repository.NewChatRepository(dbPool)

	// 4. Интеграции (Внешние API)
	llmClient := llm.NewClient(cfg.LLMBaseURL, cfg.LLMApiKey)
	rerankClient := rerank.NewRerankClient(cfg.LLMBaseURL, cfg.LLMApiKey)

	// 5. Сервисы (Бизнес-логика)
	authSvc := service.NewAuthService(userRepo, cfg)
	bankSvc := service.NewBankService(accRepo, txRepo)
	aiAnalysisSvc := service.NewAnalysisService(*llmClient, accRepo, finRepo)
	aiAssistantSvc := service.NewAssistantService(*llmClient, rerankClient, ragRepo, chatRepo)

	// 6. API (Роутер и Хэндлеры)
	router := api.SetupRouter(authSvc, bankSvc, aiAnalysisSvc, aiAssistantSvc)

	// 7. Запуск
	log.Printf("Zaman AI Assistant starting on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
