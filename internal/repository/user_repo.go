package repository

import (
	"ai-assistant/internal/domain"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	CreateUser(ctx context.Context, email, passHash, fullName string) (*domain.User, error)
	FindUserByEmail(ctx context.Context, email string) (*domain.User, error)
	FindUserByID(ctx context.Context, id int) (*domain.User, error)
}
type pgUserRepo struct { db *pgxpool.Pool }
func NewUserRepository(db *pgxpool.Pool) UserRepository { return &pgUserRepo{db} }

func (r *pgUserRepo) CreateUser(ctx context.Context, email, passHash, fullName string) (*domain.User, error) {
	u := &domain.User{}
	query := `INSERT INTO users (email, password_hash, full_name) 
               VALUES ($1, $2, $3) 
               RETURNING id, email, full_name, created_at`
	err := r.db.QueryRow(ctx, query, email, passHash, fullName).Scan(&u.ID, &u.Email, &u.FullName, &u.CreatedAt)
	return u, err
}

func (r *pgUserRepo) FindUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	u := &domain.User{}
	query := `SELECT id, email, password_hash, full_name, created_at FROM users WHERE email = $1`
	err := r.db.QueryRow(ctx, query, email).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.FullName, &u.CreatedAt)
	return u, err
}

func (r *pgUserRepo) FindUserByID(ctx context.Context, id int) (*domain.User, error) {
	u := &domain.User{}
	query := `SELECT id, email, password_hash, full_name, created_at FROM users WHERE id = $1`
	err := r.db.QueryRow(ctx, query, id).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.FullName, &u.CreatedAt)
	return u, err
}