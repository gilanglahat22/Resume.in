package controllers

import (
	"net/http"
	"strings"
	"time"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
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
	resumeRepo  models.ResumeRepository
}

// NewChatbotController creates a new chatbot controller
func NewChatbotController(chatbotRepo models.ChatbotRepository, resumeRepo models.ResumeRepository) *ChatbotController {
	return &ChatbotController{
		chatbotRepo: chatbotRepo,
		resumeRepo:  resumeRepo,
	}
}

// SendMessage handles the chat API endpoint
// @Summary Send a message to the chatbot
// @Description Send a message to the chatbot and get a response
// @Tags chatbot
// @Accept json
// @Produce json
// @Param request body ChatRequest true "Chat request"
// @Success 200 {object} docs.ChatResponse
// @Failure 400 {object} docs.ErrorResponse
// @Failure 500 {object} docs.ErrorResponse
// @Router /chat/message [post]
func (c *ChatbotController) SendMessage(ctx *gin.Context) {
	var request ChatRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// If session ID is not provided, create a new one
	if request.SessionID == "" {
		request.SessionID = utils.GenerateUUID()
		utils.Info("Created new session ID: %s", request.SessionID)
	}
	
	utils.Info("Processing chat message for session ID: %s", request.SessionID)

	// Process the query
	response, err := c.chatbotRepo.ProcessQuery(ctx.Request.Context(), request.SessionID, request.Query)
	if err != nil {
		utils.Error("Failed to process query: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process query"})
		return
	}

	// Check if this has information relevant for resume generation and add a hint
	resumeKeywords := []string{
		"resume", "CV", "experience", "job", "work", "skill", "education",
		"qualification", "degree", "university", "college", "project",
	}
	
	isResumeRelated := false
	for _, keyword := range resumeKeywords {
		if strings.Contains(strings.ToLower(request.Query), strings.ToLower(keyword)) {
			isResumeRelated = true
			break
		}
	}
	
	responseData := gin.H{
		"session_id": request.SessionID,
		"response":   response,
	}
	
	if isResumeRelated {
		responseData["resume_hint"] = true
		responseData["resume_message"] = "I've saved this information for your resume. When you're ready, you can generate your resume by sending a request to the generate-resume endpoint."
	}

	// Return the response
	ctx.JSON(http.StatusOK, responseData)
}

// GetChatHistory retrieves the chat history for a session
// @Summary Get chat history
// @Description Get the chat history for a specific session
// @Tags chatbot
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID"
// @Success 200 {object} docs.ChatHistoryResponse
// @Failure 400 {object} docs.ErrorResponse
// @Failure 500 {object} docs.ErrorResponse
// @Router /chat/history/{sessionId} [get]
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
// @Summary Upload a document
// @Description Upload a document to the vector store for context retrieval
// @Tags chatbot
// @Accept json
// @Produce json
// @Param request body UploadDocumentRequest true "Document upload request"
// @Success 200 {object} docs.DocumentResponse
// @Failure 400 {object} docs.ErrorResponse
// @Failure 500 {object} docs.ErrorResponse
// @Router /chat/document [post]
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

// GenerateResumeRequest contains the information needed to generate an ATS resume
type GenerateResumeRequest struct {
	SessionID string `json:"session_id" binding:"required"`
	Query     string `json:"query" binding:"omitempty"` // Make query optional for compatibility with chat requests
}

// GenerateATSResume generates an ATS-optimized resume in PDF format from chat data
// @Summary Generate ATS Resume
// @Description Process chat history to generate an ATS-formatted resume in PDF
// @Tags chatbot
// @Accept json
// @Produce application/pdf
// @Param request body GenerateResumeRequest true "Generate resume request"
// @Success 200 {file} file "Resume PDF file"
// @Failure 400 {object} docs.ErrorResponse
// @Failure 500 {object} docs.ErrorResponse
// @Router /chat/generate-resume [post]
func (c *ChatbotController) GenerateATSResume(ctx *gin.Context) {
	var request GenerateResumeRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// Validate the sessionID is provided
	if request.SessionID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "session_id is required"})
		return
	}

	utils.Info("Generating resume for session ID: %s", request.SessionID)
	
	// If a query is provided, process it first to add it to the chat history
	if request.Query != "" {
		utils.Info("Processing query before generating resume: %s", request.Query)
		
		// Create an embedding for the query (simplified for now)
		embedding := make([]float32, 1536)
		for i := range embedding {
			embedding[i] = float32(i % 100) / 100.0
		}
		
		// Save the user message
		userMsg := models.ChatMessage{
			SessionID: request.SessionID,
			Role:      "user",
			Content:   request.Query,
			Embedding: embedding,
			CreatedAt: time.Now(),
		}
		
		_, saveErr := c.chatbotRepo.SaveMessage(ctx.Request.Context(), userMsg)
		if saveErr != nil {
			utils.Warning("Failed to save user query to chat history: %v", saveErr)
			// Continue anyway, we still have the session ID
		} else {
			utils.Info("Saved query to chat history")
			
			// Process the query to get a response from chatbot
			_, procErr := c.chatbotRepo.ProcessQuery(ctx.Request.Context(), request.SessionID, request.Query)
			if procErr != nil {
				utils.Warning("Failed to get chatbot response for query: %v", procErr)
				// Continue anyway, we don't need the response for resume generation
			} else {
				utils.Info("Processed query and saved assistant response")
			}
		}
	}

	// Get chat history from the session (including any newly added messages)
	messages, err := c.chatbotRepo.GetSessionMessages(ctx.Request.Context(), request.SessionID)
	if err != nil {
		utils.Error("Failed to get chat history: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get chat history"})
		return
	}

	// Check if there are messages
	if len(messages) == 0 {
		utils.Warning("No messages found for session %s", request.SessionID)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No chat history found for this session"})
		return
	}

	utils.Info("Found %d messages for resume generation", len(messages))

	// Extract resume data from chat messages
	resumeData, err := c.extractResumeDataFromChat(messages)
	if err != nil {
		utils.Error("Failed to extract resume data: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to extract resume data from chat"})
		return
	}

	// Generate PDF file
	pdfPath, err := c.generatePDF(resumeData)
	if err != nil {
		utils.Error("Failed to generate PDF: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate PDF"})
		return
	}

	utils.Info("Resume PDF generated successfully at: %s", pdfPath)

	// Set filename for download
	filename := "ATS_Resume_" + resumeData.BasicInfo.Name + ".pdf"
	filename = strings.ReplaceAll(filename, " ", "_")

	// Serve the file
	ctx.FileAttachment(pdfPath, filename)

	// Clean up the temporary file after serving
	go func() {
		time.Sleep(5 * time.Second)
		os.Remove(pdfPath)
	}()
}

// extractResumeDataFromChat processes chat messages to extract structured resume data
func (c *ChatbotController) extractResumeDataFromChat(messages []models.ChatMessage) (models.Resume, error) {
	// Initialize an empty resume
	resume := models.Resume{
		ID:         utils.GenerateUUID(),
		BasicInfo:  models.BasicInfo{},
		Experience: []models.Experience{},
		Education:  []models.Education{},
		Skills:     []models.Skill{},
		Projects:   []models.Project{},
	}

	// Process all chat messages to extract information
	var allUserContent string
	for _, msg := range messages {
		if msg.Role == "user" {
			allUserContent += msg.Content + " "
			
			// Extract basic info
			if strings.Contains(strings.ToLower(msg.Content), "name is") || strings.Contains(strings.ToLower(msg.Content), "my name") {
				parts := strings.Split(msg.Content, "is")
				if len(parts) > 1 {
					namePart := strings.TrimSpace(parts[1])
					endIdx := strings.IndexAny(namePart, ".!?")
					if endIdx > 0 {
						namePart = namePart[:endIdx]
					}
					resume.BasicInfo.Name = strings.TrimSpace(namePart)
				}
			}
			
			// Extract email
			if strings.Contains(strings.ToLower(msg.Content), "email") || strings.Contains(strings.ToLower(msg.Content), "@") {
				words := strings.Fields(msg.Content)
				for _, word := range words {
					if strings.Contains(word, "@") && strings.Contains(word, ".") {
						resume.BasicInfo.Email = strings.Trim(word, ".,;:!?")
					}
				}
			}
			
			// Extract phone
			if strings.Contains(strings.ToLower(msg.Content), "phone") || strings.Contains(strings.ToLower(msg.Content), "contact") {
				words := strings.Fields(msg.Content)
				for i, word := range words {
					if strings.Contains(strings.ToLower(word), "phone") && i < len(words)-1 {
						resume.BasicInfo.Phone = strings.Trim(words[i+1], ".,;:!?")
					}
				}
			}
			
			// Extract skills
			if strings.Contains(strings.ToLower(msg.Content), "skill") || strings.Contains(strings.ToLower(msg.Content), "know") || strings.Contains(strings.ToLower(msg.Content), "can do") {
				skillKeywords := []string{"programming", "language", "framework", "software", "tool", "technology"}
				
				for _, keyword := range skillKeywords {
					if strings.Contains(strings.ToLower(msg.Content), keyword) {
						parts := strings.Split(strings.ToLower(msg.Content), keyword)
						if len(parts) > 1 {
							skillsPart := parts[1]
							endIdx := strings.IndexAny(skillsPart, ".!?")
							if endIdx > 0 {
								skillsPart = skillsPart[:endIdx]
							}
							
							skillsList := strings.Split(skillsPart, ",")
							for _, skill := range skillsList {
								skill = strings.TrimSpace(skill)
								if skill != "" && !strings.Contains(strings.ToLower(skill), "and") {
									// Check if skill already exists
									exists := false
									for _, existingSkill := range resume.Skills {
										if existingSkill.Name == skill {
											exists = true
											break
										}
									}
									
									if !exists {
										resume.Skills = append(resume.Skills, models.Skill{
											Name:     strings.TrimSpace(skill),
											Category: keyword,
										})
									}
								}
							}
						}
					}
				}
			}
			
			// Extract experience
			if strings.Contains(strings.ToLower(msg.Content), "work") || strings.Contains(strings.ToLower(msg.Content), "job") || strings.Contains(strings.ToLower(msg.Content), "experience") {
				lines := strings.Split(msg.Content, ".")
				for _, line := range lines {
					if strings.Contains(strings.ToLower(line), "work") || strings.Contains(strings.ToLower(line), "job") || strings.Contains(strings.ToLower(line), "experience") {
						exp := models.Experience{
							Description: strings.TrimSpace(line) + ".",
						}
						
						// Try to extract company name
						companies := []string{"at", "for", "with"}
						for _, company := range companies {
							if strings.Contains(strings.ToLower(line), company+" ") {
								parts := strings.Split(strings.ToLower(line), company+" ")
								if len(parts) > 1 {
									companyName := strings.Split(parts[1], " ")[0]
									exp.Company = strings.TrimSpace(companyName)
									break
								}
							}
						}
						
						// Try to extract position
						positions := []string{"as a", "as an", "position", "role"}
						for _, position := range positions {
							if strings.Contains(strings.ToLower(line), position+" ") {
								parts := strings.Split(strings.ToLower(line), position+" ")
								if len(parts) > 1 {
									positionName := strings.Split(parts[1], " ")[0]
									exp.Position = strings.TrimSpace(positionName)
									break
								}
							}
						}
						
						// Try to extract dates
						months := []string{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December", "Jan", "Feb", "Mar", "Apr", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
						years := []string{"2020", "2021", "2022", "2023", "2024", "2025"}
						
						for _, month := range months {
							if strings.Contains(line, month) {
								for _, year := range years {
									if strings.Contains(line, year) {
										exp.StartDate = month + " " + year
										break
									}
								}
							}
						}
						
						if exp.StartDate == "" {
							exp.StartDate = "2020-01" // Default
						}
						
						if strings.Contains(strings.ToLower(line), "present") || strings.Contains(strings.ToLower(line), "current") || strings.Contains(strings.ToLower(line), "now") {
							exp.EndDate = "Present"
						} else {
							exp.EndDate = "Present" // Default to present
						}
						
						// Check if we have the minimum required info
						if exp.Description != "" && (exp.Company != "" || exp.Position != "") {
							if exp.Company == "" {
								exp.Company = "Unknown Company"
							}
							if exp.Position == "" {
								exp.Position = "Software Developer" // Default
							}
							resume.Experience = append(resume.Experience, exp)
						}
					}
				}
			}
		}
	}
	
	// Fill in defaults if data couldn't be extracted
	if resume.BasicInfo.Name == "" {
		resume.BasicInfo.Name = "Job Applicant"
	}
	
	if resume.BasicInfo.Email == "" {
		resume.BasicInfo.Email = "applicant@example.com"
	}
	
	// Generate a summary from all user messages
	if len(allUserContent) > 0 {
		summary := "Professional with experience in "
		
		// Add skills to summary
		if len(resume.Skills) > 0 {
			for i, skill := range resume.Skills {
				if i > 0 && i < len(resume.Skills)-1 {
					summary += ", "
				} else if i > 0 {
					summary += " and "
				}
				summary += skill.Name
			}
		} else {
			summary += "software development and technology"
		}
		
		// Add positions to summary
		if len(resume.Experience) > 0 {
			summary += ". Previous roles include "
			for i, exp := range resume.Experience {
				if i > 0 && i < len(resume.Experience)-1 {
					summary += ", "
				} else if i > 0 {
					summary += " and "
				}
				summary += exp.Position
				if exp.Company != "" {
					summary += " at " + exp.Company
				}
			}
		}
		
		summary += "."
		resume.Summary = summary
	} else {
		resume.Summary = "Professional summary extracted from chat conversation."
	}
	
	// Add a default experience if none was extracted
	if len(resume.Experience) == 0 {
		resume.Experience = append(resume.Experience, models.Experience{
			Company:     "Example Company",
			Position:    "Software Developer",
			StartDate:   "2020-01",
			EndDate:     "Present",
			Description: "Worked on various software development projects.",
		})
	}
	
	// Add default education if none was extracted
	if len(resume.Education) == 0 {
		resume.Education = append(resume.Education, models.Education{
			Institution: "University Example",
			Degree:      "Bachelor's",
			Field:       "Computer Science",
			StartDate:   "2016-09",
			EndDate:     "2020-05",
		})
	}
	
	// Add some default skills if none were extracted
	if len(resume.Skills) == 0 {
		defaultSkills := []string{"Programming", "Problem Solving", "Communication", "Teamwork"}
		for _, skill := range defaultSkills {
			resume.Skills = append(resume.Skills, models.Skill{
				Name:     skill,
				Category: "General",
			})
		}
	}

	return resume, nil
}

// generatePDF creates an ATS-optimized PDF resume
func (c *ChatbotController) generatePDF(resume models.Resume) (string, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	
	// Set up fonts
	pdf.SetFont("Helvetica", "B", 16)
	
	// Name and contact info
	pdf.Cell(190, 10, resume.BasicInfo.Name)
	pdf.Ln(12)
	
	pdf.SetFont("Helvetica", "", 10)
	if resume.BasicInfo.Email != "" {
		pdf.Cell(190, 6, "Email: "+resume.BasicInfo.Email)
		pdf.Ln(6)
	}
	if resume.BasicInfo.Phone != "" {
		pdf.Cell(190, 6, "Phone: "+resume.BasicInfo.Phone)
		pdf.Ln(6)
	}
	if resume.BasicInfo.LinkedIn != "" {
		pdf.Cell(190, 6, "LinkedIn: "+resume.BasicInfo.LinkedIn)
		pdf.Ln(6)
	}
	
	// Summary
	pdf.Ln(4)
	pdf.SetFont("Helvetica", "B", 12)
	pdf.Cell(190, 8, "PROFESSIONAL SUMMARY")
	pdf.Ln(8)
	pdf.SetFont("Helvetica", "", 10)
	pdf.MultiCell(190, 5, resume.Summary, "", "", false)
	
	// Experience
	pdf.Ln(4)
	pdf.SetFont("Helvetica", "B", 12)
	pdf.Cell(190, 8, "EXPERIENCE")
	pdf.Ln(8)
	
	for _, exp := range resume.Experience {
		pdf.SetFont("Helvetica", "B", 10)
		pdf.Cell(190, 6, exp.Position+" | "+exp.Company)
		pdf.Ln(6)
		
		pdf.SetFont("Helvetica", "I", 10)
		pdf.Cell(190, 6, exp.StartDate+" - "+exp.EndDate)
		pdf.Ln(6)
		
		pdf.SetFont("Helvetica", "", 10)
		pdf.MultiCell(190, 5, exp.Description, "", "", false)
		
		if len(exp.Highlights) > 0 {
			pdf.Ln(2)
			for _, highlight := range exp.Highlights {
				pdf.Cell(5, 5, "•")
				pdf.Cell(185, 5, highlight)
				pdf.Ln(5)
			}
		}
		
		pdf.Ln(4)
	}
	
	// Education
	if len(resume.Education) > 0 {
		pdf.SetFont("Helvetica", "B", 12)
		pdf.Cell(190, 8, "EDUCATION")
		pdf.Ln(8)
		
		for _, edu := range resume.Education {
			pdf.SetFont("Helvetica", "B", 10)
			pdf.Cell(190, 6, edu.Degree+" in "+edu.Field)
			pdf.Ln(6)
			
			pdf.SetFont("Helvetica", "", 10)
			pdf.Cell(190, 6, edu.Institution)
			pdf.Ln(6)
			
			pdf.SetFont("Helvetica", "I", 10)
			pdf.Cell(190, 6, edu.StartDate+" - "+edu.EndDate)
			pdf.Ln(8)
		}
	}
	
	// Skills
	if len(resume.Skills) > 0 {
		pdf.Ln(4)
		pdf.SetFont("Helvetica", "B", 12)
		pdf.Cell(190, 8, "SKILLS")
		pdf.Ln(8)
		
		pdf.SetFont("Helvetica", "", 10)
		var skillText string
		for i, skill := range resume.Skills {
			skillText += skill.Name
			if i < len(resume.Skills)-1 {
				skillText += " • "
			}
		}
		pdf.MultiCell(190, 5, skillText, "", "", false)
	}
	
	// Create output directory if it doesn't exist
	outputDir := filepath.Join("test", "resume_pdfs")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", err
	}
	
	// Save PDF to test folder
	outputPath := filepath.Join(outputDir, "resume_"+utils.GenerateUUID()+".pdf")
	err := pdf.OutputFileAndClose(outputPath)
	if err != nil {
		return "", err
	}
	
	return outputPath, nil
} 