package service

import (
	"ai-assistant/internal/domain"
	"ai-assistant/internal/repository"
	"context"
	"errors"
)
type BankService interface {
	MakePayment(ctx context.Context, userID int, amount float64, merchant, category string) (*domain.Transaction, error)
	//... (другие функции, e.g., TopUp)
}
type bankSvc struct {
	accRepo repository.AccountRepository
	txRepo  repository.TransactionRepository
}
func NewBankService(accRepo repository.AccountRepository, txRepo repository.TransactionRepository) BankService {
	return &bankSvc{accRepo: accRepo, txRepo: txRepo}
}

func (s *bankSvc) MakePayment(ctx context.Context, userID int, amount float64, merchant, category string) (*domain.Transaction, error) {
	acc, err := s.accRepo.GetAccountByUserID(ctx, userID)
	if err != nil { return nil, err }

	if acc.Balance < amount { return nil, errors.New("insufficient funds") }

	newBalance := acc.Balance - amount
	if err := s.accRepo.UpdateBalance(ctx, userID, newBalance); err != nil {
		return nil, err
	}

	tx := &domain.Transaction{
		AccountID: acc.ID,
		Amount:    amount,
		Merchant:  merchant,
		Category:  category,
		Type:      "expense",
	}
	return s.txRepo.CreateTransaction(ctx, tx)
}