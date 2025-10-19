package api

import (
	"ai-assistant/internal/service"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// (Этот файл без изменений, просто для контекста)
// type Handlers struct {
// 	Auth *AuthHandlers
// 	Bank *BankHandlers
// 	Chat *ChatHandlers
// }

func SetupRouter(
	authSvc service.AuthService,
	bankSvc service.BankService,
	aiAnalysisSvc service.AnalysisService,
	aiAssistantSvc service.AssistantService,
) *chi.Mux {

	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"}, // (Настроить для прода)
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type"},
	}))

	// --- Инициализация Хэндлеров ---
	authHandlers := NewAuthHandlers(authSvc) // (!!!) Создаем хэндлеры
	bankHandlers := NewBankHandlers(bankSvc, aiAnalysisSvc)
	chatHandlers := NewChatHandlers(aiAssistantSvc)

	// --- Роуты ---

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// /api/v1/auth/... (Публичные роуты)
	r.Route("/api/v1/auth", func(r chi.Router) {
		r.Post("/register", authHandlers.HandleRegister) // (!!!) Добавлен
		r.Post("/login", authHandlers.HandleLogin)       // (!!!) Добавлен
	})

	// /api/v1/app/... (Защищенные роуты)
	r.Route("/api/v1/app", func(r chi.Router) {
		r.Use(AuthMiddleware(authSvc)) // (!!!) Защита

		// Банкинг (с проактивным AI)
		r.Post("/payment", bankHandlers.HandlePayment)

		// Чат (с реактивным AI)
		r.Post("/chat", chatHandlers.HandleChat)
		// r.Get("/chat/stream", chatHandlers.HandleChatStream) // (WebSocket)
		// r.Post("/chat/voice", chatHandlers.HandleVoice)      // (Whisper)
	})

	return r
}
