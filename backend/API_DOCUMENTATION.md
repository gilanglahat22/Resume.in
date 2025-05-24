# Resume.in API Documentation

## Swagger UI
Access the interactive API documentation at: http://localhost:8080/swagger/index.html

## Authentication
Most endpoints require JWT authentication. Use the `Authorization` header with the Bearer scheme:
```
Authorization: Bearer {your-jwt-token}
```

## API Endpoints

### Authentication Endpoints

#### 1. Google OAuth Login
- **GET** `/api/auth/google/login`
- **Authentication**: Not required
- **Description**: Initiates Google OAuth flow and returns the authorization URL
- **Response**: 
  ```json
  {
    "auth_url": "https://accounts.google.com/o/oauth2/v2/auth?..."
  }
  ```

#### 2. Google OAuth Callback
- **GET** `/api/auth/google/callback`
- **Authentication**: Not required
- **Query Parameters**: 
  - `code` (required): Authorization code from Google
  - `state` (required): OAuth state for security
- **Description**: Processes the OAuth callback and returns JWT tokens
- **Response**: 
  ```json
  {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "email": "user@example.com",
      "name": "John Doe",
      "picture": "https://example.com/profile.jpg",
      "provider": "google",
      "role": "user",
      "created_at": "2023-05-17T01:52:36.789Z",
      "updated_at": "2023-05-17T01:52:36.789Z"
    },
    "expires_in": 86400
  }
  ```

#### 3. Refresh Token
- **POST** `/api/auth/refresh`
- **Authentication**: Not required
- **Request Body**:
  ```json
  {
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
  ```
- **Description**: Refresh an expired JWT token using refresh token
- **Response**: Same as Google OAuth Callback

#### 4. Logout
- **POST** `/api/auth/logout`
- **Authentication**: Required (Bearer token)
- **Description**: Logout the current user
- **Response**: 
  ```json
  {
    "message": "Logged out successfully"
  }
  ```

#### 5. Get User Profile
- **GET** `/api/auth/profile`
- **Authentication**: Required (Bearer token)
- **Description**: Get the current authenticated user's profile
- **Response**: User object

### Resume Endpoints

#### 1. Get All Resumes
- **GET** `/api/resumes`
- **Authentication**: Required (Bearer token)
- **Description**: Get a list of all available resumes
- **Response**: Array of Resume objects

#### 2. Get Resume by ID
- **GET** `/api/resumes/{id}`
- **Authentication**: Required (Bearer token)
- **Path Parameters**: 
  - `id` (required): Resume ID
- **Description**: Get a specific resume by its ID
- **Response**: Resume object

#### 3. Create Resume
- **POST** `/api/resumes`
- **Authentication**: Required (Bearer token)
- **Request Body**: Resume object
- **Description**: Add a new resume to the system
- **Response**: Created Resume object

#### 4. Update Resume
- **PUT** `/api/resumes/{id}`
- **Authentication**: Required (Bearer token)
- **Path Parameters**: 
  - `id` (required): Resume ID
- **Request Body**: Resume object
- **Description**: Update an existing resume by its ID
- **Response**: Updated Resume object

#### 5. Delete Resume
- **DELETE** `/api/resumes/{id}`
- **Authentication**: Required (Bearer token)
- **Path Parameters**: 
  - `id` (required): Resume ID
- **Description**: Delete a resume by its ID
- **Response**: 
  ```json
  {
    "status": "deleted"
  }
  ```

### Chatbot Endpoints

#### 1. Send Message
- **POST** `/api/chat/message`
- **Authentication**: Required (Bearer token)
- **Request Body**:
  ```json
  {
    "query": "What can you tell me about resume formatting?",
    "session_id": "user123" // optional
  }
  ```
- **Description**: Send a message to the chatbot and get a response
- **Response**: 
  ```json
  {
    "session_id": "user123",
    "response": {
      "answer": "A well-formatted resume should be clean...",
      "sources": ["resume-guide.pdf", "formatting-tips.txt"],
      "created_at": "2023-05-17T01:52:36.789Z"
    },
    "resume_hint": true, // optional
    "resume_message": "I've saved this information for your resume..." // optional
  }
  ```

#### 2. Get Chat History
- **GET** `/api/chat/history/{sessionId}`
- **Authentication**: Required (Bearer token)
- **Path Parameters**: 
  - `sessionId` (required): Session ID
- **Description**: Get the chat history for a specific session
- **Response**: 
  ```json
  {
    "session_id": "user123",
    "messages": [
      {
        "id": 1,
        "session_id": "user123",
        "role": "user",
        "content": "What can you tell me about resume formatting?",
        "created_at": "2023-05-17T01:52:36.789Z"
      }
    ]
  }
  ```

#### 3. Upload Document
- **POST** `/api/chat/document`
- **Authentication**: Required (Bearer token)
- **Request Body**:
  ```json
  {
    "content": "A resume is a document that summarizes...",
    "metadata": {
      "source": "resume-guide.pdf",
      "title": "Resume Formatting Guide"
    }
  }
  ```
- **Description**: Upload a document to the vector store for context retrieval
- **Response**: 
  ```json
  {
    "status": "success",
    "message": "Document uploaded successfully"
  }
  ```

#### 4. Generate Resume
- **POST** `/api/chat/generate-resume`
- **Authentication**: Required (Bearer token)
- **Request Body**:
  ```json
  {
    "session_id": "user123",
    "query": "Generate my resume based on our conversation" // optional
  }
  ```
- **Description**: Process chat history to generate an ATS-formatted resume in PDF
- **Response**: PDF file download

### Other Endpoints

#### 1. Health Check
- **GET** `/health`
- **Authentication**: Not required
- **Description**: Check if the API is running
- **Response**: 
  ```json
  {
    "status": "ok"
  }
  ```

## Error Responses

All endpoints may return error responses in the following format:
```json
{
  "error": "Error message description"
}
```

Common HTTP status codes:
- `400 Bad Request`: Invalid request data
- `401 Unauthorized`: Authentication required or invalid token
- `404 Not Found`: Resource not found
- `500 Internal Server Error`: Server error

## Models

### User
```json
{
  "id": "string",
  "email": "string",
  "name": "string",
  "picture": "string",
  "provider": "string",
  "role": "string",
  "created_at": "datetime",
  "updated_at": "datetime"
}
```

### Resume
```json
{
  "id": "string",
  "basic_info": {
    "name": "string",
    "email": "string",
    "phone": "string",
    "address": "string",
    "website": "string",
    "linkedin": "string",
    "github": "string"
  },
  "summary": "string",
  "experience": [
    {
      "company": "string",
      "position": "string",
      "start_date": "string",
      "end_date": "string",
      "description": "string",
      "highlights": ["string"]
    }
  ],
  "education": [
    {
      "institution": "string",
      "degree": "string",
      "field": "string",
      "start_date": "string",
      "end_date": "string",
      "gpa": "string"
    }
  ],
  "skills": [
    {
      "name": "string",
      "level": "string",
      "category": "string"
    }
  ],
  "projects": [
    {
      "name": "string",
      "description": "string",
      "start_date": "string",
      "end_date": "string",
      "url": "string",
      "technologies": ["string"]
    }
  ]
}
```

## Testing the API

### 1. Get OAuth URL
```bash
curl http://localhost:8080/api/auth/google/login
```

### 2. After OAuth login, use the token
```bash
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" http://localhost:8080/api/resumes
```

### 3. Send a chat message
```bash
curl -X POST http://localhost:8080/api/chat/message \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "query": "I am a software engineer with 5 years of experience"
  }'
```

### 4. Generate resume from chat
```bash
curl -X POST http://localhost:8080/api/chat/generate-resume \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "YOUR_SESSION_ID"
  }' \
  --output resume.pdf
``` 