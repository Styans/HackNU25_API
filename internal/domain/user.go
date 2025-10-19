package domain
import "time"
type User struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // Не отдаем в JSON
	FullName     string    `json:"full_name"`
	CreatedAt    time.Time `json:"created_at"`
}