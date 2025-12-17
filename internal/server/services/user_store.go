package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/OvsienkoValeriya/GophKeeper/internal/models"
	"github.com/OvsienkoValeriya/GophKeeper/migrations"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)

type UserStore interface {
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	GetUserByID(ctx context.Context, id int64) (*models.User, error)

	SetMasterKey(ctx context.Context, userID int64, salt, verifier []byte) error
	GetMasterKeyData(ctx context.Context, userID int64) (salt, verifier []byte, err error)
	HasMasterKey(ctx context.Context, userID int64) (bool, error)
}

type PostgresUserStore struct {
	db *sqlx.DB
}

func NewPostgresUserStore(dsn string) (*PostgresUserStore, error) {
	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	store := &PostgresUserStore{db: db}
	if err := migrations.Run(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return store, nil
}

func (s *PostgresUserStore) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	query := `
		INSERT INTO users (login, password_hash)
		VALUES ($1, $2)
		RETURNING id
	`
	err := s.db.QueryRowxContext(ctx, query, user.Username, user.Password).Scan(&user.ID)
	if err != nil {
		if isUniqueViolation(err) {
			return nil, ErrUserAlreadyExists
		}
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return user, nil
}

func isUniqueViolation(err error) bool {
	return err != nil && (err.Error() == "pq: duplicate key value violates unique constraint" ||
		err.Error() != "" && err.Error()[0:5] == "ERROR")
}

func (s *PostgresUserStore) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	query := `SELECT id, login, password_hash FROM users WHERE login = $1`
	var user models.User
	err := s.db.QueryRowContext(ctx, query, username).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

func (s *PostgresUserStore) GetUserByID(ctx context.Context, id int64) (*models.User, error) {
	query := `SELECT id, login, password_hash FROM users WHERE id = $1`
	var user models.User
	err := s.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

func (s *PostgresUserStore) Close() error {
	return s.db.Close()
}

func (s *PostgresUserStore) SetMasterKey(ctx context.Context, userID int64, salt, verifier []byte) error {
	query :=
		`UPDATE users 
	 SET master_key_salt = $1,
	  master_key_verifier = $2,
	  master_key_created_at = $3
	  WHERE id = $4`

	_, err := s.db.ExecContext(ctx, query, salt, verifier, time.Now(), userID)
	if err != nil {
		return fmt.Errorf("failed to set master key: %w", err)
	}
	return nil
}

func (s *PostgresUserStore) GetMasterKeyData(ctx context.Context, userID int64) ([]byte, []byte, error) {
	query := `SELECT master_key_salt,master_key_verifier 
	 FROM users 
	 WHERE id = $1`
	var salt, verifier []byte
	err := s.db.QueryRowContext(ctx, query, userID).Scan(&salt, &verifier)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, ErrUserNotFound
		}
		return nil, nil, fmt.Errorf("failed to get master key data: %w", err)
	}
	return salt, verifier, nil
}

func (s *PostgresUserStore) HasMasterKey(ctx context.Context, userID int64) (bool, error) {
	query := `
        SELECT master_key_salt IS NOT NULL 
        FROM users 
        WHERE id = $1
    `
	var hasMasterKey bool
	err := s.db.QueryRowContext(ctx, query, userID).Scan(&hasMasterKey)
	if err != nil {
		return false, fmt.Errorf("failed to check master key: %w", err)
	}
	return hasMasterKey, nil
}
