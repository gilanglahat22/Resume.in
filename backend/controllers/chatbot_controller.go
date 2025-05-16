package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"resume.in/backend/models"
	"resume.in/backend/utils"
)

// ChatRequest represents a chat request from the user
type ChatRequest struct {
	Query     string `json:"query" binding:"required"`
	SessionID string `json:"session_id"`
}

// ChatbotController handles chatbot-related API endpoints
type ChatbotController struct {
	chatbotRepo models.ChatbotRepository
}

// NewChatbotController creates a new chatbot controller
func NewChatbotController(chatbotRepo models.ChatbotRepository) *ChatbotController {
	return &ChatbotController{
		chatbotRepo: chatbotRepo,
	}
}

// SendMessage handles the chat API endpoint
func (c *ChatbotController) SendMessage(ctx *gin.Context) {
	var request ChatRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// If session ID is not provided, create a new one
	if request.SessionID == "" {
		request.SessionID = utils.GenerateUUID()
	}

	// Process the query
	response, err := c.chatbotRepo.ProcessQuery(ctx.Request.Context(), request.SessionID, request.Query)
	if err != nil {
		utils.Error("Failed to process query: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process query"})
		return
	}

	// Return the response
	ctx.JSON(http.StatusOK, gin.H{
		"session_id": request.SessionID,
		"response":   response,
	})
}

// GetChatHistory retrieves the chat history for a session
func (c *ChatbotController) GetChatHistory(ctx *gin.Context) {
	sessionID := ctx.Param("sessionId")
	if sessionID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Session ID is required"})
		return
	}

	// Get chat history
	messages, err := c.chatbotRepo.GetSessionMessages(ctx.Request.Context(), sessionID)
	if err != nil {
		utils.Error("Failed to get chat history: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get chat history"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"session_id": sessionID,
		"messages":   messages,
	})
}

// UploadDocumentRequest represents a document upload request
type UploadDocumentRequest struct {
	Content  string                 `json:"content" binding:"required"`
	Metadata map[string]interface{} `json:"metadata"`
}

// UploadDocument handles uploading a document to the vector store
func (c *ChatbotController) UploadDocument(ctx *gin.Context) {
	var request UploadDocumentRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// Use the AddDocument method from the interface
	if err := c.chatbotRepo.AddDocument(ctx.Request.Context(), request.Content, request.Metadata); err != nil {
		utils.Error("Failed to upload document: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload document"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"message": "Document uploaded successfully",
	})
} 