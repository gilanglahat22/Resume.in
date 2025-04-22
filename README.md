# Resume.in

A web application for creating and managing professional resumes, built with Go (Gin) for the backend, Angular for the frontend, and PostgreSQL for the database.

## Prerequisites

- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)

## Getting Started

### Using Scripts

#### On Linux/Mac:

1. Make the scripts executable:
   ```bash
   chmod +x start.sh stop.sh
   ```

2. Start the application:
   ```bash
   ./start.sh
   ```

3. Stop the application:
   ```bash
   ./stop.sh
   ```

#### On Windows:

1. Start the application:
   ```
   start.bat
   ```

2. Stop the application:
   ```
   stop.bat
   ```

### Manual Commands

1. Build and start the containers:
   ```bash
   docker-compose up --build
   ```

2. Stop the containers:
   ```bash
   docker-compose down
   ```

## Accessing the Application

- Frontend: http://localhost:4200
- Backend API: http://localhost:8080/api
- Health Check: http://localhost:8080/health

## API Endpoints

- `GET /api/resume`: Get all resumes
- `GET /api/resume/:id`: Get a resume by ID
- `POST /api/resume`: Create a new resume
- `PUT /api/resume/:id`: Update a resume
- `DELETE /api/resume/:id`: Delete a resume
- `GET /api/skills`: Get all skills from all resumes
- `GET /api/experience`: Get all experiences from all resumes

## Project Structure

```
resume.in/
├── backend/              # Go backend
│   ├── config/           # Configuration
│   ├── controllers/      # HTTP handlers
│   ├── middleware/       # Middleware
│   ├── models/           # Data models
│   ├── routes/           # API routes
│   └── utils/            # Utilities
├── client/               # Angular frontend
│   ├── src/              # Source code
│   └── ...
├── docker-compose.yml    # Docker Compose configuration
├── start.sh              # Script to start the application (Linux/Mac)
├── stop.sh               # Script to stop the application (Linux/Mac)
├── start.bat             # Script to start the application (Windows)
└── stop.bat              # Script to stop the application (Windows)
```

## Development

### Backend

To run the backend locally for development:

```bash
cd backend
go run main.go
```

### Frontend

To run the frontend locally for development:

```bash
cd client
npm install
npm start
```