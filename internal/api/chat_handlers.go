package api

import (
	"ai-assistant/internal/service"
	"encoding/json"
	"net/http"
	// "github.com/google/uuid"
	// "github.com/gorilla/websocket"
)

// var upgrader = websocket.Upgrader{
// 	CheckOrigin: func(r *http.Request) bool { return true },
// }

type ChatHandlers struct {
	aiSvc service.AssistantService
}
func NewChatHandlers(aiSvc service.AssistantService) *ChatHandlers {
	return &ChatHandlers{aiSvc: aiSvc}
}

type ChatRequest struct {
	Query     string `json:"query"`
	SessionID string `json:"session_id"` // (Клиент должен управлять сессией)
}
type ChatResponse struct {
	Response string `json:"response"`
}

// HandleChat (простой POST, без стриминга)
func (h *ChatHandlers) HandleChat(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(int) // (Из middleware)
	
	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// (В идеале SessionID надо генерировать, если пустой)
	// if req.SessionID == "" { req.SessionID = uuid.NewString() }

	response, err := h.aiSvc.GetResponse(r.Context(), userID, req.SessionID, req.Query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ChatResponse{Response: response})
}

// HandleChatStream (для WebSocket)
// func (h *ChatHandlers) HandleChatStream(w http.ResponseWriter, r *http.Request) {
// 	conn, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil { log.Println(err); return }
// 	defer conn.Close()
//  
//   // ... (Логика чтения/записи в WebSocket) ...
// 	// h.aiSvc.GetResponseStream(ctx, userID, sessionID, query, conn)
// }