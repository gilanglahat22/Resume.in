package models

import (
	"errors"
	"sync"
)

// ResumeRepository defines the interface for resume data operations
type ResumeRepository interface {
	FindAll() []Resume
	FindByID(id string) (Resume, error)
	Create(resume Resume) (Resume, error)
	Update(id string, resume Resume) (Resume, error)
	Delete(id string) error
	GetAllSkills() []Skill
	GetAllExperience() []Experience
}

// InMemoryResumeRepository implements ResumeRepository with an in-memory map
type InMemoryResumeRepository struct {
	resumes map[string]Resume
	mutex   sync.RWMutex
}

// NewInMemoryResumeRepository creates a new in-memory resume repository
func NewInMemoryResumeRepository() *InMemoryResumeRepository {
	return &InMemoryResumeRepository{
		resumes: make(map[string]Resume),
	}
}

// FindAll returns all resumes
func (r *InMemoryResumeRepository) FindAll() []Resume {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var result []Resume
	for _, resume := range r.resumes {
		result = append(result, resume)
	}
	return result
}

// FindByID returns a resume by its ID
func (r *InMemoryResumeRepository) FindByID(id string) (Resume, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	resume, exists := r.resumes[id]
	if !exists {
		return Resume{}, errors.New("resume not found")
	}
	return resume, nil
}

// Create adds a new resume
func (r *InMemoryResumeRepository) Create(resume Resume) (Resume, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if resume.ID == "" {
		return Resume{}, errors.New("resume ID is required")
	}

	if _, exists := r.resumes[resume.ID]; exists {
		return Resume{}, errors.New("resume with this ID already exists")
	}

	r.resumes[resume.ID] = resume
	return resume, nil
}

// Update modifies an existing resume
func (r *InMemoryResumeRepository) Update(id string, resume Resume) (Resume, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.resumes[id]; !exists {
		return Resume{}, errors.New("resume not found")
	}

	resume.ID = id
	r.resumes[id] = resume
	return resume, nil
}

// Delete removes a resume
func (r *InMemoryResumeRepository) Delete(id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.resumes[id]; !exists {
		return errors.New("resume not found")
	}

	delete(r.resumes, id)
	return nil
}

// GetAllSkills returns all skills from all resumes
func (r *InMemoryResumeRepository) GetAllSkills() []Skill {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var allSkills []Skill
	for _, resume := range r.resumes {
		allSkills = append(allSkills, resume.Skills...)
	}
	return allSkills
}

// GetAllExperience returns all experiences from all resumes
func (r *InMemoryResumeRepository) GetAllExperience() []Experience {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var allExperience []Experience
	for _, resume := range r.resumes {
		allExperience = append(allExperience, resume.Experience...)
	}
	return allExperience
}

// InitDemoData adds sample data to the repository
func (r *InMemoryResumeRepository) InitDemoData() {
	sampleResume := Resume{
		ID: "sample",
		BasicInfo: BasicInfo{
			Name:     "John Doe",
			Email:    "john.doe@example.com",
			Phone:    "+1 (123) 456-7890",
			Address:  "123 Main St, City, Country",
			Website:  "https://johndoe.com",
			LinkedIn: "https://linkedin.com/in/johndoe",
			GitHub:   "https://github.com/johndoe",
		},
		Summary: "Software engineer with 5 years of experience in web development",
		Experience: []Experience{
			{
				Company:     "Tech Company",
				Position:    "Senior Developer",
				StartDate:   "2021-01-01",
				EndDate:     "Present",
				Description: "Full-stack development with React and Go",
				Highlights: []string{
					"Improved application performance by 30%",
					"Led a team of 5 developers",
				},
			},
		},
		Education: []Education{
			{
				Institution: "University of Example",
				Degree:      "Bachelor's",
				Field:       "Computer Science",
				StartDate:   "2014-09-01",
				EndDate:     "2018-06-30",
				GPA:         "3.8",
			},
		},
		Skills: []Skill{
			{Name: "Go", Level: "Expert", Category: "Programming Languages"},
			{Name: "React", Level: "Advanced", Category: "Frontend"},
			{Name: "Docker", Level: "Intermediate", Category: "DevOps"},
		},
		Projects: []Project{
			{
				Name:        "Resume Builder",
				Description: "A web application to create and manage resumes",
				StartDate:   "2022-03-01",
				EndDate:     "2022-06-01",
				URL:         "https://github.com/johndoe/resume-builder",
				Technologies: []string{"Go", "React", "PostgreSQL"},
			},
		},
	}

	r.Create(sampleResume)
} 