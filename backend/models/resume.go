package models

// Resume represents the resume data structure
type Resume struct {
	ID           string       `json:"id"`
	BasicInfo    BasicInfo    `json:"basicInfo"`
	Summary      string       `json:"summary"`
	Experience   []Experience `json:"experience"`
	Education    []Education  `json:"education"`
	Skills       []Skill      `json:"skills"`
	Certificates []Certificate `json:"certificates"`
	Projects     []Project    `json:"projects"`
}

// BasicInfo contains personal details
type BasicInfo struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
	Website  string `json:"website"`
	LinkedIn string `json:"linkedin"`
	GitHub   string `json:"github"`
}

// Experience represents work experience
type Experience struct {
	Company     string   `json:"company"`
	Position    string   `json:"position"`
	StartDate   string   `json:"startDate"`
	EndDate     string   `json:"endDate"`
	Description string   `json:"description"`
	Highlights  []string `json:"highlights"`
}

// Education represents educational background
type Education struct {
	Institution string `json:"institution"`
	Degree      string `json:"degree"`
	Field       string `json:"field"`
	StartDate   string `json:"startDate"`
	EndDate     string `json:"endDate"`
	GPA         string `json:"gpa"`
}

// Skill represents a professional skill
type Skill struct {
	Name     string `json:"name"`
	Level    string `json:"level"`
	Category string `json:"category"`
}

// Certificate represents professional certifications
type Certificate struct {
	Name       string `json:"name"`
	Issuer     string `json:"issuer"`
	IssueDate  string `json:"issueDate"`
	ExpiryDate string `json:"expiryDate"`
	URL        string `json:"url"`
}

// Project represents personal or professional projects
type Project struct {
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	StartDate    string   `json:"startDate"`
	EndDate      string   `json:"endDate"`
	URL          string   `json:"url"`
	Technologies []string `json:"technologies"`
} 