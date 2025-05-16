package models

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/pgvector/pgvector-go"
	"resume.in/backend/utils"
	"resume.in/backend/config"
)

// SimplePostgresChatbotRepository provides a simplified implementation
// without LangChain dependencies to avoid build issues
type SimplePostgresChatbotRepository struct {
	db *sql.DB
}

// OpenRouterRequest represents a request to the Open Router API
type OpenRouterRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens"`
	Temperature float64   `json:"temperature,omitempty"`
}

// Message represents a chat message in the Open Router API format
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OpenRouterResponse represents a response from the Open Router API
type OpenRouterResponse struct {
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// NewPostgresChatbotRepository creates a new PostgreSQL chatbot repository
// This is the function called from main.go
func NewPostgresChatbotRepository(db *sql.DB) (ChatbotRepository, error) {
	return NewSimplePostgresChatbotRepository(db)
}

// NewSimplePostgresChatbotRepository creates a new PostgreSQL chatbot repository
func NewSimplePostgresChatbotRepository(db *sql.DB) (*SimplePostgresChatbotRepository, error) {
	repo := &SimplePostgresChatbotRepository{
		db: db,
	}
	
	// Initialize tables
	if err := repo.initTables(); err != nil {
		return nil, err
	}
	
	return repo, nil
}

// initTables creates the necessary tables for the chatbot functionality
func (r *SimplePostgresChatbotRepository) initTables() error {
	// Create extension if not exists
	_, err := r.db.Exec(`CREATE EXTENSION IF NOT EXISTS vector`)
	if err != nil {
		return err
	}
	
	// Create chat_messages table
	_, err = r.db.Exec(`
		CREATE TABLE IF NOT EXISTS chat_messages (
			id SERIAL PRIMARY KEY,
			session_id VARCHAR(255) NOT NULL,
			role VARCHAR(50) NOT NULL,
			content TEXT NOT NULL,
			embedding vector(1536),
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)
	`)
	if err != nil {
		return err
	}
	
	// Create vector_documents table
	_, err = r.db.Exec(`
		CREATE TABLE IF NOT EXISTS vector_documents (
			id VARCHAR(255) PRIMARY KEY,
			content TEXT NOT NULL,
			metadata JSONB,
			embedding vector(1536) NOT NULL
		)
	`)
	if err != nil {
		return err
	}
	
	// Create index on the embedding column
	_, err = r.db.Exec(`
		CREATE INDEX IF NOT EXISTS vector_documents_embedding_idx 
		ON vector_documents 
		USING ivfflat (embedding vector_cosine_ops)
		WITH (lists = 100)
	`)
	
	return err
}

// SaveMessage saves a chat message to the database
func (r *SimplePostgresChatbotRepository) SaveMessage(ctx context.Context, message ChatMessage) (ChatMessage, error) {
	query := `
		INSERT INTO chat_messages (session_id, role, content, embedding, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	var id int64
	var embedVector pgvector.Vector

	if message.Embedding != nil {
		embedVector = pgvector.NewVector(message.Embedding)
	}

	err := r.db.QueryRowContext(
		ctx,
		query,
		message.SessionID,
		message.Role,
		message.Content,
		embedVector,
		time.Now(),
	).Scan(&id)

	if err != nil {
		return ChatMessage{}, err
	}

	message.ID = id
	return message, nil
}

// GetSessionMessages retrieves all messages for a given session
func (r *SimplePostgresChatbotRepository) GetSessionMessages(ctx context.Context, sessionID string) ([]ChatMessage, error) {
	query := `
		SELECT id, session_id, role, content, created_at
		FROM chat_messages
		WHERE session_id = $1
		ORDER BY created_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []ChatMessage
	for rows.Next() {
		var message ChatMessage
		if err := rows.Scan(
			&message.ID,
			&message.SessionID,
			&message.Role,
			&message.Content,
			&message.CreatedAt,
		); err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	return messages, nil
}

// StoreDocument stores a document in the vector database
func (r *SimplePostgresChatbotRepository) StoreDocument(ctx context.Context, doc VectorDocument) error {
	query := `
		INSERT INTO vector_documents (id, content, metadata, embedding)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) 
		DO UPDATE SET content = $2, metadata = $3, embedding = $4
	`

	if doc.ID == "" {
		doc.ID = uuid.New().String()
	}

	metadata, err := json.Marshal(doc.Metadata)
	if err != nil {
		return err
	}

	// Use a simple float32 array with values between 0-1 for demo purposes
	// In a real implementation, this would be generated by an embedding model
	if doc.Embedding == nil {
		doc.Embedding = make([]float32, 1536)
		for i := range doc.Embedding {
			doc.Embedding[i] = float32(i % 100) / 100.0
		}
	}

	embedVector := pgvector.NewVector(doc.Embedding)
	_, err = r.db.ExecContext(
		ctx,
		query,
		doc.ID,
		doc.Content,
		metadata,
		embedVector,
	)

	return err
}

// SearchSimilarDocuments searches for similar documents in the vector database
func (r *SimplePostgresChatbotRepository) SearchSimilarDocuments(ctx context.Context, embedding []float32, limit int) ([]VectorDocument, error) {
	query := `
		SELECT id, content, metadata
		FROM vector_documents
		ORDER BY embedding <=> $1
		LIMIT $2
	`

	if limit <= 0 {
		limit = 5
	}

	embedVector := pgvector.NewVector(embedding)
	rows, err := r.db.QueryContext(ctx, query, embedVector, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var documents []VectorDocument
	for rows.Next() {
		var doc VectorDocument
		var metadataBytes []byte

		if err := rows.Scan(&doc.ID, &doc.Content, &metadataBytes); err != nil {
			return nil, err
		}

		if len(metadataBytes) > 0 {
			if err := json.Unmarshal(metadataBytes, &doc.Metadata); err != nil {
				return nil, err
			}
		}

		documents = append(documents, doc)
	}

	return documents, nil
}

// callOpenRouter sends a request to Open Router API and returns the response
func callOpenRouter(messages []Message) (string, error) {
	// Get configuration for Open Router
	cfg := config.LoadConfigFromEnv()
	if cfg.OpenRouterAPIKey == "" || cfg.OpenRouterAPIKey == "your_openrouter_api_key" {
		utils.Error("OPEN_ROUTER_API_KEY is not set or is using the default value. Please set a valid API key.")
		return "I'm sorry, but my connection to the language model is not configured correctly. Please check your OPEN_ROUTER_API_KEY environment variable.", nil
	}

	// Prepare request to Open Router API
	requestData := OpenRouterRequest{
		Model:       cfg.OpenRouterModel, // Default to a model specified in config
		Messages:    messages,
		MaxTokens:   1000, // Lower token limit to ensure it stays within free tier
		Temperature: 0.7,  // Add temperature for more balanced responses
	}

	// Log model being used
	utils.Info("Using Open Router model: %s with max_tokens: %d", cfg.OpenRouterModel, requestData.MaxTokens)

	requestBody, err := json.Marshal(requestData)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+cfg.OpenRouterAPIKey)
	req.Header.Set("HTTP-Referer", "https://resume.in") // Replace with your actual domain
	req.Header.Set("X-Title", "Resume.in Chatbot")

	// Send request
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	// Parse response
	var openRouterResp OpenRouterResponse
	if err := json.Unmarshal(body, &openRouterResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	// Check for errors
	if openRouterResp.Error != nil {
		return "", fmt.Errorf("open router error: %s", openRouterResp.Error.Message)
	}

	// Check if we have valid choices
	if len(openRouterResp.Choices) == 0 {
		return "", fmt.Errorf("no response from LLM")
	}

	// Return the content of the first choice
	return openRouterResp.Choices[0].Message.Content, nil
}

// ProcessQuery processes a user query and returns a response using Open Router
func (r *SimplePostgresChatbotRepository) ProcessQuery(ctx context.Context, sessionID string, query string) (ChatResponse, error) {
	// Create a simple embedding (in a real implementation, this would use an embedding model)
	embedding := make([]float32, 1536)
	for i := range embedding {
		embedding[i] = float32(i % 100) / 100.0
	}

	// Save the user message
	userMsg := ChatMessage{
		SessionID: sessionID,
		Role:      "user",
		Content:   query,
		Embedding: embedding,
		CreatedAt: time.Now(),
	}
	
	_, err := r.SaveMessage(ctx, userMsg)
	if err != nil {
		utils.Error("Failed to save user message: %v", err)
		// Continue processing anyway
	}

	// Search for similar documents
	docs, err := r.SearchSimilarDocuments(ctx, embedding, 5)
	if err != nil {
		return ChatResponse{}, fmt.Errorf("failed to search similar documents: %w", err)
	}

	// Prepare context for LLM from retrieved documents
	var context string
	var sources []string
	
	if len(docs) > 0 {
		context = "Context information:\n\n"
		for _, doc := range docs {
			context += doc.Content + "\n\n"
			if doc.Metadata != nil {
				if source, ok := doc.Metadata["source"].(string); ok {
					sources = append(sources, source)
				}
			}
		}
	}

	// Get previous conversation history
	history, err := r.GetSessionMessages(ctx, sessionID)
	if err != nil {
		utils.Warning("Failed to get conversation history: %v", err)
		// Continue with empty history
	}

	// Prepare messages for Open Router
	var messages []Message
	
	// System message with context
	messages = append(messages, Message{
		Role:    "system",
		Content: "You are a helpful assistant. Use the following context information to answer the user's question, if relevant: " + context,
	})

	// Add conversation history (up to 10 messages to avoid token limits)
	maxHistory := 10
	if len(history) > maxHistory {
		history = history[len(history)-maxHistory:]
	}
	
	for _, msg := range history {
		messages = append(messages, Message{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	// Add the current query if not already in history
	if len(history) == 0 || history[len(history)-1].Content != query {
		messages = append(messages, Message{
			Role:    "user",
			Content: query,
		})
	}

	// Call Open Router API
	response, err := callOpenRouter(messages)
	if err != nil {
		utils.Error("Failed to get response from Open Router: %v", err)
		response = "I'm sorry, I'm having trouble processing your request right now. Please try again later."
	}

	// Save the assistant message - make sure to include the embedding
	assistantMsg := ChatMessage{
		SessionID: sessionID,
		Role:      "assistant",
		Content:   response,
		Embedding: embedding, // Use the same embedding as the query for simplicity
		CreatedAt: time.Now(),
	}
	
	_, err = r.SaveMessage(ctx, assistantMsg)
	if err != nil {
		utils.Error("Failed to save assistant message: %v", err)
		// Continue anyway
	}

	// Return the response
	return ChatResponse{
		Answer:    response,
		Sources:   sources,
		CreatedAt: time.Now(),
	}, nil
}

// AddDocument is a helper function that prepares and stores a document
func (r *SimplePostgresChatbotRepository) AddDocument(ctx context.Context, content string, metadata map[string]interface{}) error {
	// Create a simple embedding (in a real implementation, this would use an embedding model)
	embedding := make([]float32, 1536)
	for i := range embedding {
		embedding[i] = float32(i % 100) / 100.0
	}

	doc := VectorDocument{
		ID:        uuid.New().String(),
		Content:   content,
		Metadata:  metadata,
		Embedding: embedding,
	}

	return r.StoreDocument(ctx, doc)
}

// Helper function for min
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
} 