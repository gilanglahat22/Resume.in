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
func (c *ResumeController) GetResume(ctx *gin.Context) {
	id := ctx.Param("id")
	
	resume, err := c.repository.FindByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	
	ctx.JSON(http.StatusOK, resume)
}

// GetAllResumes retrieves all resumes
func (c *ResumeController) GetAllResumes(ctx *gin.Context) {
	resumes := c.repository.FindAll()
	ctx.JSON(http.StatusOK, resumes)
}

// CreateResume adds a new resume
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
func (c *ResumeController) GetAllSkills(ctx *gin.Context) {
	skills := c.repository.GetAllSkills()
	ctx.JSON(http.StatusOK, skills)
}

// GetAllExperience retrieves all experiences from all resumes
func (c *ResumeController) GetAllExperience(ctx *gin.Context) {
	experiences := c.repository.GetAllExperience()
	ctx.JSON(http.StatusOK, experiences)
} 