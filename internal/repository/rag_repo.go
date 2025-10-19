package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pgvector/pgvector-go"
)

type RAGRepository interface {
	FindRelevantDocs(ctx context.Context, embedding []float32, limit int) ([]string, error)
}
type pgRAGRepo struct{ db *pgxpool.Pool }

func NewRAGRepository(db *pgxpool.Pool) RAGRepository { return &pgRAGRepo{db} }

func (r *pgRAGRepo) FindRelevantDocs(ctx context.Context, embedding []float32, limit int) ([]string, error) {
	var docs []string
	query := `SELECT content FROM product_embeddings ORDER BY embedding <-> $1 LIMIT $2`
	rows, err := r.db.Query(ctx, query, pgvector.NewVector(embedding), limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var content string
		if err := rows.Scan(&content); err != nil {
			return nil, err
		}
		docs = append(docs, content)
	}
	return docs, nil
}
