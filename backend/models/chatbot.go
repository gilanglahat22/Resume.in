package models

import (
	"context"
	"time"
)

// ChatMessage represents a message in a chat conversation
type ChatMessage struct {
	ID        int64     `json:"id"`
	SessionID string    `json:"session_id"`
	Role      string    `json:"role"` // "user" or "assistant"
	Content   string    `json:"content"`
	Embedding []float32 `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

// ChatQuery represents a user query and its embedding
type ChatQuery struct {
	Query     string    `json:"query"`
	Embedding []float32 `json:"-"`
}

// ChatResponse represents a response from the LLM
type ChatResponse struct {
	Answer    string    `json:"answer"`
	Sources   []string  `json:"sources,omitempty"`
	Metadata  any       `json:"metadata,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// VectorDocument represents a document stored in the vector database
type VectorDocument struct {
	ID        string                 `json:"id"`
	Content   string                 `json:"content"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	Embedding []float32              `json:"-"`
}

// ChatbotRepository defines the interface for chatbot operations
type ChatbotRepository interface {
	// Message management
	SaveMessage(ctx context.Context, message ChatMessage) (ChatMessage, error)
	GetSessionMessages(ctx context.Context, sessionID string) ([]ChatMessage, error)
	
	// Vector operations
	StoreDocument(ctx context.Context, doc VectorDocument) error
	SearchSimilarDocuments(ctx context.Context, embedding []float32, limit int) ([]VectorDocument, error)
	
	// Query handling
	ProcessQuery(ctx context.Context, sessionID string, query string) (ChatResponse, error)
	
	// Document management helper
	AddDocument(ctx context.Context, content string, metadata map[string]interface{}) error
} 