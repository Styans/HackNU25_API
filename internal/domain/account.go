package domain
import "time"
type Account struct {
	ID             int       `json:"id"`
	UserID         int       `json:"user_id"`
	CardNumberMock string    `json:"card_number_mock"`
	Balance        float64   `json:"balance"`
	NextPayday     time.Time `json:"next_payday"`
}