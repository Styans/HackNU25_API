package api

import (
	"ai-assistant/internal/service"
	"encoding/json"
	"log"
	"net/http"
)

type BankHandlers struct {
	bankSvc service.BankService
	aiSvc   service.AnalysisService // AI-анализатор
}

func NewBankHandlers(bankSvc service.BankService, aiSvc service.AnalysisService) *BankHandlers {
	return &BankHandlers{bankSvc: bankSvc, aiSvc: aiSvc}
}

type PaymentRequest struct {
	Amount   float64 `json:"amount"`
	Merchant string  `json:"merchant"`
	Category string  `json:"category"`
}
type PaymentResponse struct {
	Status   string `json:"status"`
	TxID     int    `json:"tx_id"`
	AIAdvice string `json:"ai_advice,omitempty"` // Наш проактивный совет
}

// HandlePayment - Эндпоинт для "псевдо-оплаты"
func (h *BankHandlers) HandlePayment(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(int) // (Из middleware)

	var req PaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 1. Проводим платеж
	tx, err := h.bankSvc.MakePayment(r.Context(), userID, req.Amount, req.Merchant, req.Category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity) // e.g., Insufficient funds
		return
	}

	// 2. (!!!) Вызываем AI-анализ
	advice, err := h.aiSvc.AnalyzeTransactionProactively(r.Context(), userID, tx)
	if err != nil {
		// Не страшно, если AI упал, транзакция прошла
		log.Printf("Failed to get AI advice: %v", err)
	}

	// 3. Отвечаем
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(PaymentResponse{
		Status:   "success",
		TxID:     tx.ID,
		AIAdvice: advice, // Отправляем совет клиенту!
	})
}
