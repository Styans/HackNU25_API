package repository

import (
	"ai-assistant/internal/domain"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionRepository interface {
	CreateTransaction(ctx context.Context, tx *domain.Transaction) (*domain.Transaction, error)
	GetTransactionsByAccountID(ctx context.Context, accountID int, limit int) ([]domain.Transaction, error)
}
type pgTxRepo struct { db *pgxpool.Pool }
func NewTransactionRepository(db *pgxpool.Pool) TransactionRepository { return &pgTxRepo{db} }

func (r *pgTxRepo) CreateTransaction(ctx context.Context, tx *domain.Transaction) (*domain.Transaction, error) {
	query := `INSERT INTO transactions (account_id, amount, merchant, category, type) 
               VALUES ($1, $2, $3, $4, $5) 
               RETURNING id, created_at`
	err := r.db.QueryRow(ctx, query, tx.AccountID, tx.Amount, tx.Merchant, tx.Category, tx.Type).Scan(&tx.ID, &tx.CreatedAt)
	return tx, err
}

func (r *pgTxRepo) GetTransactionsByAccountID(ctx context.Context, accountID int, limit int) ([]domain.Transaction, error) {
	// ... (Реализация) ...
	return nil, nil // Оставим для краткости
}