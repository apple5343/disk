package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	// pg driver.
	_ "github.com/jackc/pgx/v5/stdlib"
	"gitlab.crja72.ru/golang/2025/autumn/projects/go36/disk/services/user-service.v2/internal/model"
)

type PostgresRepository struct {
	db *sql.DB
}

type UserRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	GetByLogin(ctx context.Context, login string) (*model.InternalUser, error)
	Create(ctx context.Context, user *model.User) error
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

var ErrUserNotFound = errors.New("user not found")

func (r *PostgresRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	user := &model.User{}
	err := r.db.QueryRowContext(ctx,
		"SELECT id, login, name FROM userdata WHERE id = $1", id).
		Scan(&user.ID, &user.Login, &user.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func (r *PostgresRepository) GetByLogin(ctx context.Context, login string) (*model.InternalUser, error) {
	user := &model.InternalUser{}
	err := r.db.QueryRowContext(ctx,
		"SELECT id, login, name, password_hash FROM userdata WHERE login = $1", login).
		Scan(&user.ID, &user.Login, &user.Name, &user.PasswordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func (r *PostgresRepository) Create(ctx context.Context, user *model.User) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO userdata (id, login, name, password_hash) VALUES ($1, $2, $3, $4)",
		user.ID, user.Login, user.Name, user.PasswordHash,
	)
	return err
}
