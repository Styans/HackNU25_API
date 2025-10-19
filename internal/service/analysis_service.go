package service

import (
	"ai-assistant/internal/domain"
	"ai-assistant/internal/integration/llm"
	"ai-assistant/internal/repository"
	"context"
	"fmt"
	"time"
)
type AnalysisService interface {
	AnalyzeTransactionProactively(ctx context.Context, userID int, tx *domain.Transaction) (string, error)
}
type analysisSvc struct {
	llmClient llm.LLMClient
	accRepo   repository.AccountRepository
	finRepo   repository.FinancingRepository
}
func NewAnalysisService(llm llm.LLMClient, accRepo repository.AccountRepository, finRepo repository.FinancingRepository) AnalysisService {
	return &analysisSvc{llmClient: llm, accRepo: accRepo, finRepo: finRepo}
}

func (s *analysisSvc) AnalyzeTransactionProactively(ctx context.Context, userID int, tx *domain.Transaction) (string, error) {
	acc, err := s.accRepo.GetAccountByUserID(ctx, userID)
	if err != nil { return "", err }
	
	fins, err := s.finRepo.GetFinancingByUserID(ctx, userID)
	if err != nil { return "", err }

	daysToPayday := int(time.Until(acc.NextPayday).Hours() / 24)

	// (Логика из нашего прошлого обсуждения)
	isLowBalance := acc.Balance < 50000 // Порог
	isFarFromPayday := daysToPayday > 14

	if !isLowBalance || !isFarFromPayday {
		return "", nil // Нет совета, все ОК
	}

	// Формируем промпт
	prompt := fmt.Sprintf(`
        Ты - "Zaman Companion", заботливый финансовый коуч Halal-банка. 
        Твой клиент ТОЛЬКО ЧТО совершил покупку.

        ДАННЫЕ О ПОКУПКЕ:
        - Сумма: %.2f тг
        - Куда: "%s" (Категория: %s)

        ФИНАНСОВОЕ ПОЛОЖЕНИЕ КЛИЕНТА:
        - Остаток на счете: %.2f тг
        - Следующая зарплата: через %d дней
        - Активных финансирований: %d

        ТВОЯ ЗАДАЧА:
        Клиент в "красной зоне": у него мало денег, а до зарплаты далеко.
        Дай ОЧЕНЬ мягкий, поддерживающий совет (1-2 предложения).
        ЗАПРЕТЫ (ХАРАМ): Нельзя советовать микрозаймы, кредиты.
        
        Пример: "Заметил вашу покупку. У вас осталось %.2f тг, а до зарплаты %d дней. 
        Возможно, стоит быть немного экономнее? Давайте вместе спланируем бюджет."
        `,
		tx.Amount, tx.Merchant, tx.Category, 
		acc.Balance, daysToPayday, len(fins),
		acc.Balance, daysToPayday,
	)

	response, err := s.llmClient.GetChatCompletion(ctx, []llm.ChatMessage{{Role: "user", Content: prompt}})
	if err != nil { return "", err }

	return response, nil
}