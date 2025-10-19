package service

import (
	"ai-assistant/internal/integration/llm"
	// (!!!) ИСПРАВЛЕНИЕ: Пакет должен импортироваться БЕЗ псевдонима "service"
	"ai-assistant/internal/integration/rerank" 
	"ai-assistant/internal/repository"
	"context"
	"fmt"
	"strings"
	// "github.com/gorilla/websocket" // (Нужен для streaming)
)

type AssistantService interface {
	GetResponse(ctx context.Context, userID int, sessionID, query string) (string, error)
	// GetResponseStream(ctx context.Context, userID int, sessionID, query string, conn *websocket.Conn) error
}

type assistantSvc struct {
	llmClient    llm.LLMClient
	rerankClient rerank.RerankClient // (!!!) Теперь 'rerank' будет найдено
	ragRepo      repository.RAGRepository
	chatRepo     repository.ChatRepository
}

func NewAssistantService(llm llm.LLMClient, rerank rerank.RerankClient, ragRepo repository.RAGRepository, chatRepo repository.ChatRepository) AssistantService {
	return &assistantSvc{
		llmClient:    llm,
		rerankClient: rerank,
		ragRepo:      ragRepo,
		chatRepo:     chatRepo,
	}
}

// GetResponse - главная RAG-функция
func (s *assistantSvc) GetResponse(ctx context.Context, userID int, sessionID, query string) (string, error) {
	// 1. Модерация
	// flagged, err := s.llmClient.CheckModeration(ctx, query)
	// if err != nil {
	// 	return "", err
	// }
	// if flagged {
	// 	return "Ваш запрос не соответствует политике безопасности.", nil
	// }

	// 2. Сохранить запрос юзера
	_ = s.chatRepo.SaveMessage(ctx, userID, sessionID, "user", query)

	// // 3. RAG: Векторизация
	// embedding, err := s.llmClient.GetEmbedding(ctx, query)
	// if err != nil {
	// 	return "", err
	// }

	// 4. RAG: Поиск
	// docs, err := s.ragRepo.FindRelevantDocs(ctx, embedding, 10)
	// if err != nil {
	// 	return "", err
	// }

	// 5. RAG: Rerank (Улучшение)
	rerankedDocs, err := s.rerankClient.Rerank(ctx, query, nil)
	if err != nil {
		return "", err
	}

	// 6. Сборка Промпта
	// (Получаем историю чата)
	history, _ := s.chatRepo.GetHistoryBySessionID(ctx, sessionID, 5)

	// (!!!) ИСПРАВЛЕНИЕ: убираем неиспользуемую переменную `prompt`
	_, messages := s.buildRAGPrompt(query, rerankedDocs, history)

	// 7. Генерация ответа
	response, err := s.llmClient.GetChatCompletion(ctx, messages)
	if err != nil {
		return "", err
	}

	// 8. Сохранить ответ AI
	_ = s.chatRepo.SaveMessage(ctx, userID, sessionID, "ai", response)

	return response, nil
}

// buildRAGPrompt - "Личность" нашего AI
func (s *assistantSvc) buildRAGPrompt(query string, contextDocs []string, history []repository.ChatHistory) (string, []llm.ChatMessage) {

	contextStr := strings.Join(contextDocs, "\n---\n")

	// Системный промпт (личность)
	systemPrompt := fmt.Sprintf(`
Ты - "Zaman Companion", вежливый и умный финансовый ассистент Zaman Bank.
Ты строго следуешь принципам Исламского финансирования.
Используй ТОЛЬКО предоставленный КОНТЕКСТ для ответа на вопросы о продуктах.
Не придумывай продукты или условия.

!!! ВАЖНОЕ ПРАВИЛО (ХАЛЯЛЬ) !!!
Если клиент спрашивает про "акции", "фондовый рынок", "криптовалюту", "биткоин", "фьючерсы" или "микрозаймы":
1. Вежливо объясни, что Zaman Bank не работает с этими инструментами, так как они содержат элементы "Гарар" (неопределенность) и/или "Риба" (проценты), что не соответствует Шариату.
2. Вместо этого, предложи клиенту Халяльную альтернативу из КОНТЕКСТА, например, Исламский депозит "Аманат" (Мудараба).

КОНТЕКСТ (Продукты Zaman Bank):
---
%s
---
`, contextStr)

	var messages []llm.ChatMessage
	messages = append(messages, llm.ChatMessage{Role: "system", Content: systemPrompt})

	// Добавляем историю
	for _, msg := range history {
		messages = append(messages, llm.ChatMessage{Role: msg.Role, Content: msg.Content})
	}

	// Добавляем текущий вопрос
	messages = append(messages, llm.ChatMessage{Role: "user", Content: query})

	return systemPrompt, messages
}