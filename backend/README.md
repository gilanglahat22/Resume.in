# Resume.in Chatbot Backend

This is the backend for Resume.in application with chatbot capability powered by LLM.

## System Design

The chatbot system follows this design:

1. User sends a query to the backend
2. The query is embedded using OpenAI embeddings
3. The embedding is used to search for relevant documents in a vector database (PostgreSQL with pgvector)
4. Retrieved documents along with the query are sent to an LLM (OpenAI) for processing
5. The LLM generates a response which is sent back to the user

## Setup

### Prerequisites

- Go 1.18 or higher
- PostgreSQL with pgvector extension
- OpenAI API key

### Environment Variables

Create a `.env` file with the following variables:

```
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=resume_db

# Server Configuration
SERVER_PORT=8080

# OpenAI Configuration
OPENAI_API_KEY=your_openai_api_key
```

### Installation

#### Local Development

1. Clone the repository
2. Install dependencies:
   ```
   go mod download
   ```
3. Run the application:
   ```
   go run main.go
   ```

#### Docker Setup

1. Create a `.env` file in the root directory with:
   ```
   OPENAI_API_KEY=your_openai_api_key
   ```

2. Build and start the containers:
   ```
   docker-compose up -d
   ```

3. The services will be available at:
   - Backend: http://localhost:8080
   - Frontend: http://localhost:3000

4. To stop the containers:
   ```
   docker-compose down
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
- **LangChain (Go)**: Framework for working with LLMs
- **OpenAI**: For generating embeddings and LLM responses 