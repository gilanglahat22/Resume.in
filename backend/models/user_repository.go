package models

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
)

// PostgresUserRepository implements UserRepository using PostgreSQL
type PostgresUserRepository struct {
	db *sqlx.DB
}

// NewPostgresUserRepository creates a new PostgreSQL user repository
func NewPostgresUserRepository(db *sqlx.DB) UserRepository {
	return &PostgresUserRepository{db: db}
}

// Create creates a new user
func (r *PostgresUserRepository) Create(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (id, email, name, provider, provider_id, picture, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	
	if user.Role == "" {
		user.Role = "user"
	}
	
	_, err := r.db.ExecContext(ctx, query, 
		user.ID, 
		user.Email, 
		user.Name, 
		user.Provider, 
		user.ProviderID, 
		user.Picture, 
		user.Role,
		user.CreatedAt,
		user.UpdatedAt,
	)
	
	return err
}

// GetByID retrieves a user by ID
func (r *PostgresUserRepository) GetByID(ctx context.Context, id string) (*User, error) {
	query := `
		SELECT id, email, name, provider, provider_id, picture, role, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	
	var user User
	err := r.db.GetContext(ctx, &user, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	
	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *PostgresUserRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, email, name, provider, provider_id, picture, role, created_at, updated_at
		FROM users
		WHERE email = $1
	`
	
	var user User
	err := r.db.GetContext(ctx, &user, query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	
	return &user, nil
}

// Update updates a user
func (r *PostgresUserRepository) Update(ctx context.Context, user *User) error {
	query := `
		UPDATE users
		SET email = $2, name = $3, picture = $4, role = $5, updated_at = $6
		WHERE id = $1
	`
	
	user.UpdatedAt = time.Now()
	
	_, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.Email,
		user.Name,
		user.Picture,
		user.Role,
		user.UpdatedAt,
	)
	
	return err
}

// Delete deletes a user
func (r *PostgresUserRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE id = $1`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return errors.New("user not found")
	}
	
	return nil
}

// CreateUserTable creates the users table if it doesn't exist
func CreateUserTable(db *sqlx.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS users (
			id VARCHAR(255) PRIMARY KEY,
			email VARCHAR(255) UNIQUE NOT NULL,
			name VARCHAR(255) NOT NULL,
			provider VARCHAR(50) NOT NULL,
			provider_id VARCHAR(255) NOT NULL,
			picture TEXT,
			role VARCHAR(50) NOT NULL DEFAULT 'user',
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(provider, provider_id)
		);
		
		CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
		CREATE INDEX IF NOT EXISTS idx_users_provider_id ON users(provider, provider_id);
	`
	
	_, err := db.Exec(query)
	return err
} 