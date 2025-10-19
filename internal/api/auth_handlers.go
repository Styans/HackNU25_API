package api

import (
	"ai-assistant/internal/domain"
	"ai-assistant/internal/service"
	"encoding/json"
	"net/http"
)

type AuthHandlers struct {
	authSvc service.AuthService
}

func NewAuthHandlers(authSvc service.AuthService) *AuthHandlers {
	return &AuthHandlers{authSvc: authSvc}
}

// --- Структуры для запросов/ответов ---

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	User  *domain.User `json:"user"`
	Token string       `json:"token"`
}

// --- Хэндлеры ---

// HandleRegister обрабатывает регистрацию нового пользователя
func (h *AuthHandlers) HandleRegister(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" || req.FullName == "" {
		http.Error(w, "Email, password, and full name are required", http.StatusBadRequest)
		return
	}

	// Вызываем сервис для регистрации
	user, token, err := h.authSvc.RegisterUser(r.Context(), req.Email, req.Password, req.FullName)
	if err != nil {
		// (В проде здесь нужна более детальная обработка, e.g., "user already exists")
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	// Отправляем успешный ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(AuthResponse{
		User:  user,
		Token: token,
	})
}

// HandleLogin обрабатывает вход пользователя
func (h *AuthHandlers) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	// Вызываем сервис для логина
	user, token, err := h.authSvc.LoginUser(r.Context(), req.Email, req.Password)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Отправляем успешный ответ
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(AuthResponse{
		User:  user,
		Token: token,
	})
}