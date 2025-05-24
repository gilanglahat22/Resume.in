package models

import (
	"context"
	"time"
)

// User represents a user in the system
type User struct {
	ID          string    `json:"id" db:"id"`
	Email       string    `json:"email" db:"email"`
	Name        string    `json:"name" db:"name"`
	Password    string    `json:"-" db:"password"`  // Password is never sent to client
	Provider    string    `json:"provider" db:"provider"`     // oauth provider (google, github, etc)
	ProviderID  string    `json:"provider_id,omitempty" db:"provider_id"`        // ID from the OAuth provider
	Picture     string    `json:"picture,omitempty" db:"picture"`
	Role        string    `json:"role" db:"role"`           // user, admin
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// UserRepository defines the interface for user data access
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByProviderID(ctx context.Context, provider, providerID string) (*User, error)
	Delete(ctx context.Context, id string) error
} 