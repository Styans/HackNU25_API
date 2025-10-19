package main

import (
	"ai-assistant/internal/config"
	"ai-assistant/internal/integration/llm"
	"context"
	"log"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pgvector/pgvector-go"
)

// Сюда вы вставите ваши знания
var knowledgeBase = []string{
	"Депозит 'Аманат' - это Исламский вклад на основе принципа Мудараба. Минимальная сумма 15 000 тенге, срок от 3 до 24 месяцев. Выплата вознаграждения - ежемесячно.",
	"Мурабаха - это торговая сделка по финансированию. Банк покупает товар (например, авто или квартиру) и продает его вам в рассрочку с заранее известной наценкой. Это не кредит, так как нет процента (риба).",
	"Иджара - это Исламский лизинг. Банк покупает имущество и сдает его вам в аренду. По окончании срока аренды имущество переходит в вашу собственность.",
	"Zaman Bank не работает с криптовалютой и акциями, так как эти инструменты содержат элементы Гарар (чрезмерная неопределенность) и Риба (проценты), что не соответствует Шариату.",
	"Для накопления на Хадж или Умру мы рекомендуем депозит 'Аманат'. Он позволяет копить средства в соответствии с принципами Ислама.",
}


func main() {
	ctx := context.Background()
	log.Println("Starting RAG embedder...")
	
	cfg := config.Load()
	
	// Подключаемся к БД (используем localhost, т.к. запускаем с хост-машины)
	dbURL := "postgres://zaman_user:zaman_pass@localhost:5433/zaman_db?sslmode=disable"
	dbPool, err := pgxpool.New(ctx, dbURL)
	if err != nil { log.Fatalf("Failed to connect to DB: %v", err) }
	defer dbPool.Close()

	// Клиент к LLM
	llmClient := llm.NewClient(cfg.LLMBaseURL, cfg.LLMApiKey)
	
	log.Println("Embedding and inserting knowledge base...")

	for _, doc := range knowledgeBase {
		// 1. Получить вектор
		embedding, err := llmClient.GetEmbedding(ctx, doc)
		if err != nil {
			log.Printf("Failed to get embedding for doc: %v", err)
			continue
		}

		// 2. Вставить в БД
		query := `INSERT INTO product_embeddings (content, embedding) VALUES ($1, $2)`
		_, err = dbPool.Exec(ctx, query, doc, pgvector.NewVector(embedding))
		if err != nil {
			log.Printf("Failed to insert embedding: %v", err)
			continue
		}
		log.Printf("Successfully embedded: %s...", doc[:30])
	}
	
	log.Println("Embedding process finished.")
}