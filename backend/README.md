# Resume.in Chatbot Backend

This is the backend for Resume.in application with chatbot capability powered by LLM through Open Router.

## System Design

The chatbot system follows this design:

1. User sends a query to the backend
2. The query is processed and stored with a simple embedding
3. The embedding is used to search for relevant documents in a vector database (PostgreSQL with pgvector)
4. Retrieved documents along with the query are sent to an LLM (via Open Router) for processing
5. The LLM generates a response which is sent back to the user

## Setup

### Prerequisites

- Go 1.18 or higher
- PostgreSQL with pgvector extension
- Open Router API key (sign up at https://openrouter.ai)

### Environment Variables

Create a `.env` file with the following variables:

```
# Database Configuration
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=resumeuser
POSTGRES_PASSWORD=resumepassword
POSTGRES_DB=resumedb
POSTGRES_SSLMODE=disable

# Server Configuration
SERVER_PORT=8080
ENVIRONMENT=development
ALLOW_ORIGINS=http://localhost:3000
LOG_LEVEL=debug

# Open Router Configuration
OPEN_ROUTER_API_KEY=your_openrouter_api_key
OPEN_ROUTER_MODEL=anthropic/claude-3-sonnet:beta
```

Available Open Router models include:
- `anthropic/claude-3-opus:beta` - Highest capability Claude model
- `anthropic/claude-3-sonnet:beta` - Great balance of intelligence and speed
- `anthropic/claude-3-haiku:beta` - Fastest and most cost-effective Claude model
- `meta-llama/llama-3-70b-instruct` - Meta's most capable open model
- `mistralai/mistral-large-latest` - Mistral's most capable model

See the [Open Router models page](https://openrouter.ai/models) for a complete list of available models.

### Installation

#### Local Development

1. Clone the repository
2. Install dependencies:
   ```
   go mod download
   ```
3. Create a `.env` file in the `backend` directory with the environment variables listed above
4. Run the application:
   ```
   go run main.go
   ```

#### Docker Setup

1. Create a `.env` file in the root directory with at least these variables:
   ```
   OPEN_ROUTER_API_KEY=your_openrouter_api_key
   OPEN_ROUTER_MODEL=anthropic/claude-3-sonnet:beta
   ```

2. Build and start the containers:
   ```
   docker-compose up -d
   ```

3. The services will be available at:
   - Backend: http://localhost:8080

4. To stop the containers:
   ```
   docker-compose down
   ```

## API Documentation

The API is documented using Swagger/OpenAPI. You can access the interactive API documentation at:

```
http://localhost:8080/swagger/index.html
```

This provides a user-friendly interface for:
- Viewing all available API endpoints
- Understanding request/response formats
- Testing API endpoints directly from the browser
- Getting detailed information about API parameters and response codes

### Generating Swagger Documentation

If you make changes to the API, you'll need to regenerate the Swagger documentation:

1. Install swag CLI tool:
   ```
   go install github.com/swaggo/swag/cmd/swag@latest
   ```
   
2. Generate the docs:
   ```
   cd backend
   swag init
   ```

## API Endpoints

### Chatbot

- `POST /api/chat/message` - Send a message to the chatbot
  - Request body:
    ```json
    {
      "query": "Your question here",
      "session_id": "optional-session-id"
    }
    ```
  - Response:
    ```json
    {
      "session_id": "session-id",
      "response": {
        "answer": "LLM response",
        "sources": ["optional source references"],
        "created_at": "timestamp"
      }
    }
    ```

- `GET /api/chat/history/:sessionId` - Get chat history for a session
  - Response:
    ```json
    {
      "session_id": "session-id",
      "messages": [
        {
          "id": 1,
          "session_id": "session-id",
          "role": "user",
          "content": "User message",
          "created_at": "timestamp"
        },
        {
          "id": 2,
          "session_id": "session-id",
          "role": "assistant",
          "content": "Assistant response",
          "created_at": "timestamp"
        }
      ]
    }
    ```

- `POST /api/chat/document` - Upload a document to the vector store
  - Request body:
    ```json
    {
      "content": "Document content to be vectorized",
      "metadata": {
        "source": "optional source information",
        "title": "optional title"
      }
    }
    ```
  - Response:
    ```json
    {
      "status": "success",
      "message": "Document uploaded successfully"
    }
    ```

## Architecture

The backend uses the following components:

- **Gin**: Web framework for REST API
- **PostgreSQL with pgvector**: Vector database for storing and searching embeddings
- **Open Router**: API gateway for accessing various LLM models
- **Swagger/OpenAPI**: API documentation 