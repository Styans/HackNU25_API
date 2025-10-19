package repository

import (
	"ai-assistant/internal/domain"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FinancingRepository interface {
	GetFinancingByUserID(ctx context.Context, userID int) ([]domain.Financing, error)
}
type pgFinancingRepo struct { db *pgxpool.Pool }
func NewFinancingRepository(db *pgxpool.Pool) FinancingRepository { return &pgFinancingRepo{db} }

func (r *pgFinancingRepo) GetFinancingByUserID(ctx context.Context, userID int) ([]domain.Financing, error) {
	var financings []domain.Financing
	query := `SELECT id, user_id, product_name, total_amount, remaining_amount, monthly_payment 
               FROM financing WHERE user_id = $1`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil { return nil, err }
	defer rows.Close()

	for rows.Next() {
		var f domain.Financing
		err := rows.Scan(&f.ID, &f.UserID, &f.ProductName, &f.TotalAmount, &f.RemainingAmount, &f.MonthlyPayment)
		if err != nil { return nil, err }
		financings = append(financings, f)
	}
	return financings, nil
}