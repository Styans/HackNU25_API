package rerank

import (
	"context"
	"log"
	"net/http"
	"time"
)

type RerankClient interface {
	Rerank(ctx context.Context, query string, docs []string) ([]string, error)
}

// (!!!) ИСПРАВЛЕНИЕ: Мы определяем поля для rerank-клиента
// (по аналогии с LLMClient)
type neuralDeepReranker struct {
	httpClient *http.Client
	apiKey     string
	baseURL    string
}

// (!!!) ИСПРАВЛЕНИЕ: NewRerankClient теперь принимает конфиг,
// так как ему тоже нужен API-ключ
func NewRerankClient(baseURL, apiKey string) RerankClient {
	return &neuralDeepReranker{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		apiKey:     apiKey,
		baseURL:    baseURL,
	}
}

func (r *neuralDeepReranker) Rerank(ctx context.Context, query string, docs []string) ([]string, error) {
	// TODO: Реализовать на основе API neuraldeep.tech
	// Когда у нас будет API, мы добавим сюда свой собственный `doRequest`,
	// который будет методом (*neuralDeepReranker).
	log.Println("Rerank service is not implemented, returning original docs.")
	return docs, nil
}

// (!!!) ИСПРАВЛЕНИЕ:
// Функция doRequest(c *LLMClient) была удалена из этого файла.
// Она по ошибке была скопирована из `llm/client.go`.