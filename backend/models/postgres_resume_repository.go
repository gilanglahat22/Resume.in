package models

import (
	"encoding/json"
	"errors"
	"fmt"
	
	"github.com/jmoiron/sqlx"
)

// PostgresResumeRepository implements ResumeRepository with PostgreSQL
type PostgresResumeRepository struct {
	db *sqlx.DB
}

// NewPostgresResumeRepository creates a new repository with the given database connection
func NewPostgresResumeRepository(db *sqlx.DB) (*PostgresResumeRepository, error) {
	repo := &PostgresResumeRepository{
		db: db,
	}

	// Initialize database tables
	err := repo.initTables()
	if err != nil {
		return nil, err
	}

	return repo, nil
}

// initTables creates the necessary tables if they don't exist
func (r *PostgresResumeRepository) initTables() error {
	query := `
    CREATE TABLE IF NOT EXISTS resumes (
        id VARCHAR(100) PRIMARY KEY,
        data JSONB NOT NULL
    );
    `

	_, err := r.db.Exec(query)
	return err
}

// FindAll returns all resumes
func (r *PostgresResumeRepository) FindAll() []Resume {
	query := `SELECT data FROM resumes;`
	rows, err := r.db.Queryx(query)
	if err != nil {
		return []Resume{}
	}
	defer rows.Close()

	var resumes []Resume
	for rows.Next() {
		var data []byte
		if err := rows.Scan(&data); err != nil {
			continue
		}

		var resume Resume
		if err := json.Unmarshal(data, &resume); err != nil {
			continue
		}

		resumes = append(resumes, resume)
	}

	return resumes
}

// FindByID returns a resume by its ID
func (r *PostgresResumeRepository) FindByID(id string) (Resume, error) {
	query := `SELECT data FROM resumes WHERE id = $1;`
	var data []byte
	err := r.db.QueryRowx(query, id).Scan(&data)
	if err != nil {
		return Resume{}, errors.New("resume not found")
	}

	var resume Resume
	if err := json.Unmarshal(data, &resume); err != nil {
		return Resume{}, fmt.Errorf("failed to parse resume data: %v", err)
	}

	return resume, nil
}

// Create adds a new resume
func (r *PostgresResumeRepository) Create(resume Resume) (Resume, error) {
	if resume.ID == "" {
		return Resume{}, errors.New("resume ID is required")
	}

	// Check if resume already exists
	_, err := r.FindByID(resume.ID)
	if err == nil {
		return Resume{}, errors.New("resume with this ID already exists")
	}

	// Convert resume to JSON
	resumeJSON, err := json.Marshal(resume)
	if err != nil {
		return Resume{}, fmt.Errorf("failed to serialize resume: %v", err)
	}

	// Insert into database
	query := `INSERT INTO resumes (id, data) VALUES ($1, $2);`
	_, err = r.db.Exec(query, resume.ID, resumeJSON)
	if err != nil {
		return Resume{}, fmt.Errorf("failed to insert resume: %v", err)
	}

	return resume, nil
}

// Update modifies an existing resume
func (r *PostgresResumeRepository) Update(id string, resume Resume) (Resume, error) {
	// Check if resume exists
	_, err := r.FindByID(id)
	if err != nil {
		return Resume{}, errors.New("resume not found")
	}

	// Set ID to the path parameter value
	resume.ID = id

	// Convert resume to JSON
	resumeJSON, err := json.Marshal(resume)
	if err != nil {
		return Resume{}, fmt.Errorf("failed to serialize resume: %v", err)
	}

	// Update database
	query := `UPDATE resumes SET data = $1 WHERE id = $2;`
	_, err = r.db.Exec(query, resumeJSON, id)
	if err != nil {
		return Resume{}, fmt.Errorf("failed to update resume: %v", err)
	}

	return resume, nil
}

// Delete removes a resume
func (r *PostgresResumeRepository) Delete(id string) error {
	// Check if resume exists
	_, err := r.FindByID(id)
	if err != nil {
		return errors.New("resume not found")
	}

	// Delete from database
	query := `DELETE FROM resumes WHERE id = $1;`
	_, err = r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete resume: %v", err)
	}

	return nil
}

// GetAllSkills returns all skills from all resumes
func (r *PostgresResumeRepository) GetAllSkills() []Skill {
	var allSkills []Skill
	resumes := r.FindAll()

	for _, resume := range resumes {
		allSkills = append(allSkills, resume.Skills...)
	}

	return allSkills
}

// GetAllExperience returns all experiences from all resumes
func (r *PostgresResumeRepository) GetAllExperience() []Experience {
	var allExperience []Experience
	resumes := r.FindAll()

	for _, resume := range resumes {
		allExperience = append(allExperience, resume.Experience...)
	}

	return allExperience
}

// InitDemoData adds sample data to the repository
func (r *PostgresResumeRepository) InitDemoData() error {
	// Check if we already have data
	resumes := r.FindAll()
	if len(resumes) > 0 {
		return nil
	}

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

	_, err := r.Create(sampleResume)
	return err
} 