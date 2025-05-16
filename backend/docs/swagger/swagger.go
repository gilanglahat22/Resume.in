// Package swagger provides documentation for the API using Swagger 2.0.
package swagger

import "github.com/swaggo/swag"

//go:generate swag init --parseDependency --parseInternal -g ../../main.go -o ./

// SwaggerInfo holds exported Swagger info so clients can modify it
var SwaggerInfo = &swag.Spec{
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}

// This is the Swagger doc template - populated from the swag init command
const docTemplate = `{
    "swagger": "2.0",
    "info": {
        "description": "API Server for Resume.in application with chatbot capabilities",
        "title": "Resume.in API",
        "contact": {},
        "version": "1.0"
    },
    "host": "localhost:8080",
    "basePath": "/api",
    "paths": {
        "/chat/document": {
            "post": {
                "description": "Upload a document to the vector store for context retrieval",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "chatbot"
                ],
                "summary": "Upload a document",
                "parameters": [
                    {
                        "description": "Document upload request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/DocumentRequest"
                        }
                    }
                ],
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
                "description": "Get the chat history for a specific session",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "chatbot"
                ],
                "summary": "Get chat history",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session ID",
                        "name": "sessionId",
                        "in": "path",
                        "required": true
                    }
                ],
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
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/ErrorResponse"
                        }
                    }
                }
            }
        },
        "/chat/message": {
            "post": {
                "description": "Send a message to the chatbot and get a response",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "chatbot"
                ],
                "summary": "Send a message to the chatbot",
                "parameters": [
                    {
                        "description": "Chat request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/ChatRequest"
                        }
                    }
                ],
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
        "ChatHistoryResponse": {
            "type": "object",
            "properties": {
                "messages": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/ChatMessage"
                    }
                },
                "session_id": {
                    "type": "string",
                    "example": "user123"
                }
            }
        },
        "ChatMessage": {
            "type": "object",
            "properties": {
                "content": {
                    "type": "string",
                    "example": "What can you tell me about resume formatting?"
                },
                "created_at": {
                    "type": "string",
                    "example": "2023-05-17T01:52:36.789Z"
                },
                "id": {
                    "type": "integer",
                    "format": "int64",
                    "example": 1
                },
                "role": {
                    "type": "string",
                    "example": "user"
                },
                "session_id": {
                    "type": "string",
                    "example": "user123"
                }
            }
        },
        "ChatRequest": {
            "type": "object",
            "required": [
                "query"
            ],
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
        "ChatResponse": {
            "type": "object",
            "properties": {
                "response": {
                    "type": "object",
                    "properties": {
                        "answer": {
                            "type": "string",
                            "example": "A well-formatted resume should be clean, consistent, and easy to scan..."
                        },
                        "created_at": {
                            "type": "string",
                            "example": "2023-05-17T01:52:36.789Z"
                        },
                        "sources": {
                            "type": "array",
                            "items": {
                                "type": "string"
                            },
                            "example": [
                                "resume-guide.pdf",
                                "formatting-tips.txt"
                            ]
                        }
                    }
                },
                "session_id": {
                    "type": "string",
                    "example": "user123"
                }
            }
        },
        "DocumentRequest": {
            "type": "object",
            "required": [
                "content"
            ],
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
                "message": {
                    "type": "string",
                    "example": "Document uploaded successfully"
                },
                "status": {
                    "type": "string",
                    "example": "success"
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