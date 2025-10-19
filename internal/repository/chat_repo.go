package repository
import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)
type ChatRepository interface {
	SaveMessage(ctx context.Context, userID int, sessionID, role, content string) error
	GetHistoryBySessionID(ctx context.Context, sessionID string, limit int) ([]ChatHistory, error)
}
type ChatHistory struct { Role, Content string }
type pgChatRepo struct { db *pgxpool.Pool }
func NewChatRepository(db *pgxpool.Pool) ChatRepository { return &pgChatRepo{db} }

func (r *pgChatRepo) SaveMessage(ctx context.Context, userID int, sessionID, role, content string) error {
	query := `INSERT INTO chat_history (user_id, session_id, role, content) VALUES ($1, $2, $3, $4)`
	_, err := r.db.Exec(ctx, query, userID, sessionID, role, content)
	return err
}

func (r *pgChatRepo) GetHistoryBySessionID(ctx context.Context, sessionID string, limit int) ([]ChatHistory, error) {
	var history []ChatHistory
	query := `SELECT role, content FROM chat_history 
               WHERE session_id = $1 
               ORDER BY created_at DESC LIMIT $2`
	rows, err := r.db.Query(ctx, query, sessionID, limit)
	if err != nil { return nil, err }
	defer rows.Close()

	for rows.Next() {
		var h ChatHistory
		if err := rows.Scan(&h.Role, &h.Content); err != nil { return nil, err }
		history = append(history, h)
	}
	// Reverse to get chronological order
	for i, j := 0, len(history)-1; i < j; i, j = i+1, j-1 {
		history[i], history[j] = history[j], history[i]
	}
	return history, nil
}