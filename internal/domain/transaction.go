package domain
import "time"
type Transaction struct {
	ID        int       `json:"id"`
	AccountID int       `json:"account_id"`
	Amount    float64   `json:"amount"`
	Merchant  string    `json:"merchant"`
	Category  string    `json:"category"`
	Type      string    `json:"type"` // 'income' or 'expense'
	CreatedAt time.Time `json:"created_at"`
}