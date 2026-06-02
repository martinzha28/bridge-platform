package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/martinzha28/bridge-platform/backend/internal/models"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, username, passwordHash string) (*models.User, error) {
	u := &models.User{
		ID:           uuid.New(),
		Username:     username,
		PasswordHash: passwordHash,
		CreatedAt:    time.Now(),
	}
	_, err := r.db.Exec(ctx, `
		INSERT INTO users (id, username, password_hash, created_at)
		VALUES ($1, $2, $3, $4)`,
		u.ID, u.Username, u.PasswordHash, u.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	u := &models.User{}
	err := r.db.QueryRow(ctx, `
		SELECT id, username, password_hash, email, rating, games_played, karma,
		       about, nationality, systems, created_at
		FROM users WHERE username = $1`, username,
	).Scan(
		&u.ID, &u.Username, &u.PasswordHash, &u.Email, &u.Rating, &u.GamesPlayed, &u.Karma,
		&u.About, &u.Nationality, &u.Systems, &u.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	u := &models.User{}
	err := r.db.QueryRow(ctx, `
		SELECT id, username, password_hash, email, rating, games_played, karma,
		       about, nationality, systems, created_at
		FROM users WHERE id = $1`, id,
	).Scan(
		&u.ID, &u.Username, &u.PasswordHash, &u.Email, &u.Rating, &u.GamesPlayed, &u.Karma,
		&u.About, &u.Nationality, &u.Systems, &u.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *UserRepository) UpdateProfile(ctx context.Context, id uuid.UUID, about, nationality *string, systems []string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE users SET about = $2, nationality = $3, systems = $4
		WHERE id = $1`,
		id, about, nationality, systems,
	)
	return err
}
