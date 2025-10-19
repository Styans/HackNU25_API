package domain
type Financing struct {
	ID              int     `json:"id"`
	UserID          int     `json:"user_id"`
	ProductName     string  `json:"product_name"`
	TotalAmount     float64 `json:"total_amount"`
	RemainingAmount float64 `json:"remaining_amount"`
	MonthlyPayment  float64 `json:"monthly_payment"`
}