package repository

import (
	"ai-assistant/internal/domain"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AccountRepository interface {
	GetAccountByUserID(ctx context.Context, userID int) (*domain.Account, error)
	UpdateBalance(ctx context.Context, userID int, newBalance float64) error
}
type pgAccountRepo struct { db *pgxpool.Pool }
func NewAccountRepository(db *pgxpool.Pool) AccountRepository { return &pgAccountRepo{db} }

func (r *pgAccountRepo) GetAccountByUserID(ctx context.Context, userID int) (*domain.Account, error) {
	a := &domain.Account{}
	query := `SELECT id, user_id, card_number_mock, balance, next_payday FROM accounts WHERE user_id = $1`
	err := r.db.QueryRow(ctx, query, userID).Scan(&a.ID, &a.UserID, &a.CardNumberMock, &a.Balance, &a.NextPayday)
	return a, err
}

func (r *pgAccountRepo) UpdateBalance(ctx context.Context, userID int, newBalance float64) error {
	query := `UPDATE accounts SET balance = $1 WHERE user_id = $2`
	_, err := r.db.Exec(ctx, query, newBalance, userID)
	return err
}