package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"resume.in/backend/models"
)

// ResumeController handles resume-related HTTP requests
type ResumeController struct {
	repository models.ResumeRepository
}

// NewResumeController creates a new instance of ResumeController
func NewResumeController(repository models.ResumeRepository) *ResumeController {
	return &ResumeController{
		repository: repository,
	}
}

// GetResume retrieves a resume by ID
// @Summary Get a resume by ID
// @Description Get a specific resume by its ID
// @Tags resume
// @Accept json
// @Produce json
// @Param id path string true "Resume ID"
// @Success 200 {object} docs.Resume
// @Failure 404 {object} docs.ErrorResponse
// @Router /resume/{id} [get]
func (c *ResumeController) GetResume(ctx *gin.Context) {
	id := ctx.Param("id")
	
	resume, err := c.repository.FindByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	
	ctx.JSON(http.StatusOK, resume)
}

// GetResumes retrieves all resumes
// @Summary Get all resumes
// @Description Get a list of all available resumes
// @Tags resume
// @Accept json
// @Produce json
// @Success 200 {array} docs.Resume
// @Router /resumes [get]
func (c *ResumeController) GetResumes(ctx *gin.Context) {
	resumes := c.repository.FindAll()
	ctx.JSON(http.StatusOK, resumes)
}

// CreateResume adds a new resume
// @Summary Create a new resume
// @Description Add a new resume to the system
// @Tags resume
// @Accept json
// @Produce json
// @Param resume body models.Resume true "Resume object"
// @Success 201 {object} docs.Resume
// @Failure 400 {object} docs.ErrorResponse
// @Router /resume [post]
func (c *ResumeController) CreateResume(ctx *gin.Context) {
	var resume models.Resume
	
	if err := ctx.ShouldBindJSON(&resume); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// For demo, we'll use the name as ID if ID is not provided
	if resume.ID == "" && resume.BasicInfo.Name != "" {
		resume.ID = resume.BasicInfo.Name
	}
	
	createdResume, err := c.repository.Create(resume)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	ctx.JSON(http.StatusCreated, createdResume)
}

// UpdateResume modifies an existing resume
// @Summary Update a resume
// @Description Update an existing resume by its ID
// @Tags resume
// @Accept json
// @Produce json
// @Param id path string true "Resume ID"
// @Param resume body models.Resume true "Resume object"
// @Success 200 {object} docs.Resume
// @Failure 400 {object} docs.ErrorResponse
// @Failure 404 {object} docs.ErrorResponse
// @Router /resume/{id} [put]
func (c *ResumeController) UpdateResume(ctx *gin.Context) {
	id := ctx.Param("id")
	
	var resume models.Resume
	if err := ctx.ShouldBindJSON(&resume); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	updatedResume, err := c.repository.Update(id, resume)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	
	ctx.JSON(http.StatusOK, updatedResume)
}

// DeleteResume removes a resume
// @Summary Delete a resume
// @Description Delete a resume by its ID
// @Tags resume
// @Accept json
// @Produce json
// @Param id path string true "Resume ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} docs.ErrorResponse
// @Router /resume/{id} [delete]
func (c *ResumeController) DeleteResume(ctx *gin.Context) {
	id := ctx.Param("id")
	
	err := c.repository.Delete(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

// GetAllSkills retrieves all skills from all resumes
// @Summary Get all skills
// @Description Get a list of all skills from all resumes
// @Tags skills
// @Accept json
// @Produce json
// @Success 200 {array} docs.Skill
// @Router /skills [get]
func (c *ResumeController) GetAllSkills(ctx *gin.Context) {
	skills := c.repository.GetAllSkills()
	ctx.JSON(http.StatusOK, skills)
}

// GetAllExperience retrieves all experiences from all resumes
// @Summary Get all experiences
// @Description Get a list of all work experiences from all resumes
// @Tags experience
// @Accept json
// @Produce json
// @Success 200 {array} docs.Experience
// @Router /experience [get]
func (c *ResumeController) GetAllExperience(ctx *gin.Context) {
	experiences := c.repository.GetAllExperience()
	ctx.JSON(http.StatusOK, experiences)
} 