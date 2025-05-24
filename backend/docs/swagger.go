// Package docs implements Swagger documentation for the Resume.in API.
package docs

import (
	"github.com/swaggo/swag"
)

// DocJSON contains the Swagger API documentation in JSON format
var DocJSON = `{
	"swagger": "2.0",
	"info": {
		"description": "API Server for Resume.in application with chatbot capabilities and OAuth authentication",
		"title": "Resume.in API",
		"contact": {},
		"version": "1.0"
	},
	"host": "localhost:8080",
	"basePath": "/api",
	"schemes": ["http", "https"],
	"securityDefinitions": {
		"Bearer": {
			"type": "apiKey",
			"name": "Authorization",
			"in": "header",
			"description": "JWT Authorization header using the Bearer scheme. Example: 'Bearer {token}'"
		}
	},
	"paths": {
		"/auth/google/login": {
			"get": {
				"description": "Initiates Google OAuth flow and returns the authorization URL",
				"produces": ["application/json"],
				"tags": ["auth"],
				"summary": "Initiate Google OAuth login",
				"responses": {
					"200": {
						"description": "OK",
						"schema": {
							"type": "object",
							"properties": {
								"auth_url": {
									"type": "string",
									"example": "https://accounts.google.com/o/oauth2/v2/auth?..."
								}
							}
						}
					}
				}
			}
		},
		"/auth/google/callback": {
			"get": {
				"description": "Processes the OAuth callback and returns JWT tokens",
				"produces": ["application/json"],
				"tags": ["auth"],
				"summary": "Handle Google OAuth callback",
				"parameters": [{
					"type": "string",
					"description": "Authorization code",
					"name": "code",
					"in": "query",
					"required": true
				}, {
					"type": "string",
					"description": "OAuth state",
					"name": "state",
					"in": "query",
					"required": true
				}],
				"responses": {
					"200": {
						"description": "OK",
						"schema": {
							"$ref": "#/definitions/LoginResponse"
						}
					},
					"400": {
						"description": "Bad Request",
						"schema": {
							"$ref": "#/definitions/ErrorResponse"
						}
					},
					"500": {
						"description": "Internal Server Error",
						"schema": {
							"$ref": "#/definitions/ErrorResponse"
						}
					}
				}
			}
		},
		"/auth/refresh": {
			"post": {
				"description": "Refresh an expired JWT token using refresh token",
				"consumes": ["application/json"],
				"produces": ["application/json"],
				"tags": ["auth"],
				"summary": "Refresh JWT token",
				"parameters": [{
					"description": "Refresh token",
					"name": "refresh_token",
					"in": "body",
					"required": true,
					"schema": {
						"type": "object",
						"required": ["refresh_token"],
						"properties": {
							"refresh_token": {
								"type": "string",
								"example": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
							}
						}
					}
				}],
				"responses": {
					"200": {
						"description": "OK",
						"schema": {
							"$ref": "#/definitions/LoginResponse"
						}
					},
					"400": {
						"description": "Bad Request",
						"schema": {
							"$ref": "#/definitions/ErrorResponse"
						}
					},
					"401": {
						"description": "Unauthorized",
						"schema": {
							"$ref": "#/definitions/ErrorResponse"
						}
					}
				}
			}
		},
		"/auth/logout": {
			"post": {
				"security": [{"Bearer": []}],
				"description": "Logout the current user",
				"produces": ["application/json"],
				"tags": ["auth"],
				"summary": "Logout user",
				"responses": {
					"200": {
						"description": "OK",
						"schema": {
							"type": "object",
							"properties": {
								"message": {
									"type": "string",
									"example": "Logged out successfully"
								}
							}
						}
					}
				}
			}
		},
		"/auth/profile": {
			"get": {
				"security": [{"Bearer": []}],
				"description": "Get the current authenticated user's profile",
				"produces": ["application/json"],
				"tags": ["auth"],
				"summary": "Get user profile",
				"responses": {
					"200": {
						"description": "OK",
						"schema": {
							"$ref": "#/definitions/User"
						}
					},
					"401": {
						"description": "Unauthorized",
						"schema": {
							"$ref": "#/definitions/ErrorResponse"
						}
					}
				}
			}
		},
		"/resumes": {
			"get": {
				"security": [{"Bearer": []}],
				"description": "Get a list of all available resumes",
				"produces": ["application/json"],
				"tags": ["resume"],
				"summary": "Get all resumes",
				"responses": {
					"200": {
						"description": "OK",
						"schema": {
							"type": "array",
							"items": {
								"$ref": "#/definitions/Resume"
							}
						}
					},
					"401": {
						"description": "Unauthorized",
						"schema": {
							"$ref": "#/definitions/ErrorResponse"
						}
					}
				}
			},
			"post": {
				"security": [{"Bearer": []}],
				"description": "Add a new resume to the system",
				"consumes": ["application/json"],
				"produces": ["application/json"],
				"tags": ["resume"],
				"summary": "Create a new resume",
				"parameters": [{
					"description": "Resume object",
					"name": "resume",
					"in": "body",
					"required": true,
					"schema": {
						"$ref": "#/definitions/Resume"
					}
				}],
				"responses": {
					"201": {
						"description": "Created",
						"schema": {
							"$ref": "#/definitions/Resume"
						}
					},
					"400": {
						"description": "Bad Request",
						"schema": {
							"$ref": "#/definitions/ErrorResponse"
						}
					},
					"401": {
						"description": "Unauthorized",
						"schema": {
							"$ref": "#/definitions/ErrorResponse"
						}
					}
				}
			}
		},
		"/resumes/{id}": {
			"get": {
				"security": [{"Bearer": []}],
				"description": "Get a specific resume by its ID",
				"produces": ["application/json"],
				"tags": ["resume"],
				"summary": "Get a resume by ID",
				"parameters": [{
					"type": "string",
					"description": "Resume ID",
					"name": "id",
					"in": "path",
					"required": true
				}],
				"responses": {
					"200": {
						"description": "OK",
						"schema": {
							"$ref": "#/definitions/Resume"
						}
					},
					"401": {
						"description": "Unauthorized",
						"schema": {
							"$ref": "#/definitions/ErrorResponse"
						}
					},
					"404": {
						"description": "Not Found",
						"schema": {
							"$ref": "#/definitions/ErrorResponse"
						}
					}
				}
			},
			"put": {
				"security": [{"Bearer": []}],
				"description": "Update an existing resume by its ID",
				"consumes": ["application/json"],
				"produces": ["application/json"],
				"tags": ["resume"],
				"summary": "Update a resume",
				"parameters": [{
					"type": "string",
					"description": "Resume ID",
					"name": "id",
					"in": "path",
					"required": true
				}, {
					"description": "Resume object",
					"name": "resume",
					"in": "body",
					"required": true,
					"schema": {
						"$ref": "#/definitions/Resume"
					}
				}],
				"responses": {
					"200": {
						"description": "OK",
						"schema": {
							"$ref": "#/definitions/Resume"
						}
					},
					"400": {
						"description": "Bad Request",
						"schema": {
							"$ref": "#/definitions/ErrorResponse"
						}
					},
					"401": {
						"description": "Unauthorized",
						"schema": {
							"$ref": "#/definitions/ErrorResponse"
						}
					},
					"404": {
						"description": "Not Found",
						"schema": {
							"$ref": "#/definitions/ErrorResponse"
						}
					}
				}
			},
			"delete": {
				"security": [{"Bearer": []}],
				"description": "Delete a resume by its ID",
				"produces": ["application/json"],
				"tags": ["resume"],
				"summary": "Delete a resume",
				"parameters": [{
					"type": "string",
					"description": "Resume ID",
					"name": "id",
					"in": "path",
					"required": true
				}],
				"responses": {
					"200": {
						"description": "OK",
						"schema": {
							"type": "object",
							"properties": {
								"status": {
									"type": "string",
									"example": "deleted"
								}
							}
						}
					},
					"401": {
						"description": "Unauthorized",
						"schema": {
							"$ref": "#/definitions/ErrorResponse"
						}
					},
					"404": {
						"description": "Not Found",
						"schema": {
							"$ref": "#/definitions/ErrorResponse"
						}
					}
				}
			}
		},
		"/chat/message": {
			"post": {
				"security": [{"Bearer": []}],
				"description": "Send a message to the chatbot and get a response",
				"consumes": ["application/json"],
				"produces": ["application/json"],
				"tags": ["chatbot"],
				"summary": "Send a message to the chatbot",
				"parameters": [{
					"description": "Chat request",
					"name": "request",
					"in": "body",
					"required": true,
					"schema": {
						"$ref": "#/definitions/ChatRequest"
					}
				}],
				"responses": {
					"200": {
						"description": "OK",
						"schema": {
							"$ref": "#/definitions/ChatResponse"
						}
					},
					"400": {
						"description": "Bad Request",
						"schema": {
							"$ref": "#/definitions/ErrorResponse"
						}
					},
					"401": {
						"description": "Unauthorized",
						"schema": {
							"$ref": "#/definitions/ErrorResponse"
						}
					},
					"500": {
						"description": "Internal Server Error",
						"schema": {
							"$ref": "#/definitions/ErrorResponse"
						}
					}
				}
			}
		},
		"/chat/history/{sessionId}": {
			"get": {
				"security": [{"Bearer": []}],
				"description": "Get the chat history for a specific session",
				"produces": ["application/json"],
				"tags": ["chatbot"],
				"summary": "Get chat history",
				"parameters": [{
					"type": "string",
					"description": "Session ID",
					"name": "sessionId",
					"in": "path",
					"required": true
				}],
				"responses": {
					"200": {
						"description": "OK",
						"schema": {
							"$ref": "#/definitions/ChatHistoryResponse"
						}
					},
					"400": {
						"description": "Bad Request",
						"schema": {
							"$ref": "#/definitions/ErrorResponse"
						}
					},
					"401": {
						"description": "Unauthorized",
						"schema": {
							"$ref": "#/definitions/ErrorResponse"
						}
					},
					"500": {
						"description": "Internal Server Error",
						"schema": {
							"$ref": "#/definitions/ErrorResponse"
						}
					}
				}
			}
		},
		"/chat/document": {
			"post": {
				"security": [{"Bearer": []}],
				"description": "Upload a document to the vector store for context retrieval",
				"consumes": ["application/json"],
				"produces": ["application/json"],
				"tags": ["chatbot"],
				"summary": "Upload a document",
				"parameters": [{
					"description": "Document upload request",
					"name": "request",
					"in": "body",
					"required": true,
					"schema": {
						"$ref": "#/definitions/DocumentRequest"
					}
				}],
				"responses": {
					"200": {
						"description": "OK",
						"schema": {
							"$ref": "#/definitions/DocumentResponse"
						}
					},
					"400": {
						"description": "Bad Request",
						"schema": {
							"$ref": "#/definitions/ErrorResponse"
						}
					},
					"401": {
						"description": "Unauthorized",
						"schema": {
							"$ref": "#/definitions/ErrorResponse"
						}
					},
					"500": {
						"description": "Internal Server Error",
						"schema": {
							"$ref": "#/definitions/ErrorResponse"
						}
					}
				}
			}
		},
		"/chat/generate-resume": {
			"post": {
				"security": [{"Bearer": []}],
				"description": "Process chat history to generate an ATS-formatted resume in PDF",
				"consumes": ["application/json"],
				"produces": ["application/pdf"],
				"tags": ["chatbot"],
				"summary": "Generate ATS Resume",
				"parameters": [{
					"description": "Generate resume request",
					"name": "request",
					"in": "body",
					"required": true,
					"schema": {
						"$ref": "#/definitions/GenerateResumeRequest"
					}
				}],
				"responses": {
					"200": {
						"description": "OK",
						"schema": {
							"type": "file"
						}
					},
					"400": {
						"description": "Bad Request",
						"schema": {
							"$ref": "#/definitions/ErrorResponse"
						}
					},
					"401": {
						"description": "Unauthorized",
						"schema": {
							"$ref": "#/definitions/ErrorResponse"
						}
					},
					"500": {
						"description": "Internal Server Error",
						"schema": {
							"$ref": "#/definitions/ErrorResponse"
						}
					}
				}
			}
		}
	},
	"definitions": {
		"User": {
			"type": "object",
			"properties": {
				"id": {
					"type": "string",
					"example": "550e8400-e29b-41d4-a716-446655440000"
				},
				"email": {
					"type": "string",
					"example": "user@example.com"
				},
				"name": {
					"type": "string",
					"example": "John Doe"
				},
				"picture": {
					"type": "string",
					"example": "https://example.com/profile.jpg"
				},
				"provider": {
					"type": "string",
					"example": "google"
				},
				"role": {
					"type": "string",
					"example": "user"
				},
				"created_at": {
					"type": "string",
					"format": "date-time",
					"example": "2023-05-17T01:52:36.789Z"
				},
				"updated_at": {
					"type": "string",
					"format": "date-time",
					"example": "2023-05-17T01:52:36.789Z"
				}
			}
		},
		"LoginResponse": {
			"type": "object",
			"properties": {
				"token": {
					"type": "string",
					"example": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
				},
				"refresh_token": {
					"type": "string",
					"example": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
				},
				"user": {
					"$ref": "#/definitions/User"
				},
				"expires_in": {
					"type": "integer",
					"format": "int64",
					"example": 86400
				}
			}
		},
		"Resume": {
			"type": "object",
			"properties": {
				"id": {
					"type": "string",
					"example": "john-doe"
				},
				"basic_info": {
					"$ref": "#/definitions/BasicInfo"
				},
				"summary": {
					"type": "string",
					"example": "Experienced software engineer with expertise in full-stack development"
				},
				"experience": {
					"type": "array",
					"items": {
						"$ref": "#/definitions/Experience"
					}
				},
				"education": {
					"type": "array",
					"items": {
						"$ref": "#/definitions/Education"
					}
				},
				"skills": {
					"type": "array",
					"items": {
						"$ref": "#/definitions/Skill"
					}
				},
				"projects": {
					"type": "array",
					"items": {
						"$ref": "#/definitions/Project"
					}
				}
			}
		},
		"BasicInfo": {
			"type": "object",
			"properties": {
				"name": {
					"type": "string",
					"example": "John Doe"
				},
				"email": {
					"type": "string",
					"example": "john.doe@example.com"
				},
				"phone": {
					"type": "string",
					"example": "+1 (123) 456-7890"
				},
				"address": {
					"type": "string",
					"example": "123 Main St, City, Country"
				},
				"website": {
					"type": "string",
					"example": "https://johndoe.com"
				},
				"linkedin": {
					"type": "string",
					"example": "https://linkedin.com/in/johndoe"
				},
				"github": {
					"type": "string",
					"example": "https://github.com/johndoe"
				}
			}
		},
		"Experience": {
			"type": "object",
			"properties": {
				"company": {
					"type": "string",
					"example": "Tech Company"
				},
				"position": {
					"type": "string",
					"example": "Senior Software Engineer"
				},
				"start_date": {
					"type": "string",
					"example": "2021-01-01"
				},
				"end_date": {
					"type": "string",
					"example": "Present"
				},
				"description": {
					"type": "string",
					"example": "Led development of microservices architecture"
				},
				"highlights": {
					"type": "array",
					"items": {
						"type": "string"
					},
					"example": ["Improved performance by 30%", "Led team of 5 developers"]
				}
			}
		},
		"Education": {
			"type": "object",
			"properties": {
				"institution": {
					"type": "string",
					"example": "University of Example"
				},
				"degree": {
					"type": "string",
					"example": "Bachelor's"
				},
				"field": {
					"type": "string",
					"example": "Computer Science"
				},
				"start_date": {
					"type": "string",
					"example": "2014-09-01"
				},
				"end_date": {
					"type": "string",
					"example": "2018-06-30"
				},
				"gpa": {
					"type": "string",
					"example": "3.8"
				}
			}
		},
		"Skill": {
			"type": "object",
			"properties": {
				"name": {
					"type": "string",
					"example": "Go"
				},
				"level": {
					"type": "string",
					"example": "Expert"
				},
				"category": {
					"type": "string",
					"example": "Programming Languages"
				}
			}
		},
		"Project": {
			"type": "object",
			"properties": {
				"name": {
					"type": "string",
					"example": "Resume Builder"
				},
				"description": {
					"type": "string",
					"example": "A web application to create and manage resumes"
				},
				"start_date": {
					"type": "string",
					"example": "2022-03-01"
				},
				"end_date": {
					"type": "string",
					"example": "2022-06-01"
				},
				"url": {
					"type": "string",
					"example": "https://github.com/johndoe/resume-builder"
				},
				"technologies": {
					"type": "array",
					"items": {
						"type": "string"
					},
					"example": ["Go", "React", "PostgreSQL"]
				}
			}
		},
		"ChatRequest": {
			"type": "object",
			"required": ["query"],
			"properties": {
				"query": {
					"type": "string",
					"example": "What can you tell me about resume formatting?"
				},
				"session_id": {
					"type": "string",
					"example": "user123"
				}
			}
		},
		"GenerateResumeRequest": {
			"type": "object",
			"required": ["session_id"],
			"properties": {
				"session_id": {
					"type": "string",
					"example": "user123"
				},
				"query": {
					"type": "string",
					"example": "Generate my resume based on our conversation"
				}
			}
		},
		"ChatResponse": {
			"type": "object",
			"properties": {
				"session_id": {
					"type": "string",
					"example": "user123"
				},
				"response": {
					"type": "object",
					"properties": {
						"answer": {
							"type": "string",
							"example": "A well-formatted resume should be clean, consistent, and easy to scan..."
						},
						"sources": {
							"type": "array",
							"items": {
								"type": "string"
							},
							"example": ["resume-guide.pdf", "formatting-tips.txt"]
						},
						"created_at": {
							"type": "string",
							"format": "date-time",
							"example": "2023-05-17T01:52:36.789Z"
						}
					}
				},
				"resume_hint": {
					"type": "boolean",
					"example": true
				},
				"resume_message": {
					"type": "string",
					"example": "I've saved this information for your resume. When you're ready, you can generate your resume by sending a request to the generate-resume endpoint."
				}
			}
		},
		"ChatMessage": {
			"type": "object",
			"properties": {
				"id": {
					"type": "integer",
					"format": "int64",
					"example": 1
				},
				"session_id": {
					"type": "string",
					"example": "user123"
				},
				"role": {
					"type": "string",
					"example": "user"
				},
				"content": {
					"type": "string",
					"example": "What can you tell me about resume formatting?"
				},
				"created_at": {
					"type": "string",
					"format": "date-time",
					"example": "2023-05-17T01:52:36.789Z"
				}
			}
		},
		"ChatHistoryResponse": {
			"type": "object",
			"properties": {
				"session_id": {
					"type": "string",
					"example": "user123"
				},
				"messages": {
					"type": "array",
					"items": {
						"$ref": "#/definitions/ChatMessage"
					}
				}
			}
		},
		"DocumentRequest": {
			"type": "object",
			"required": ["content"],
			"properties": {
				"content": {
					"type": "string",
					"example": "A resume is a document that summarizes your education, work experience, skills..."
				},
				"metadata": {
					"type": "object",
					"example": {
						"source": "resume-guide.pdf",
						"title": "Resume Formatting Guide"
					}
				}
			}
		},
		"DocumentResponse": {
			"type": "object",
			"properties": {
				"status": {
					"type": "string",
					"example": "success"
				},
				"message": {
					"type": "string",
					"example": "Document uploaded successfully"
				}
			}
		},
		"ErrorResponse": {
			"type": "object",
			"properties": {
				"error": {
					"type": "string",
					"example": "Invalid request: query field is required"
				}
			}
		}
	}
}`

// SwaggerInfo holds the API information
var SwaggerInfo = swag.Spec{
	InfoInstanceName: "swagger",
	SwaggerTemplate:  DocJSON,
}

func init() {
	swag.Register("swagger", &SwaggerInfo)
} 