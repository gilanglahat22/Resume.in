# Resume.in 📄✨

<p align="center">
  <img src="https://img.shields.io/badge/status-active-success.svg" alt="Status">
  <img src="https://img.shields.io/badge/license-MIT-blue.svg" alt="License">
  <img src="https://img.shields.io/badge/PRs-welcome-brightgreen.svg" alt="PRs Welcome">
</p>

<p align="center">A modern, AI-powered web application for creating and managing professional resumes with ATS optimization.</p>

## ✨ Features

- 🤖 AI-powered resume generation based on chat conversations
- 📊 ATS optimization to help your resume pass applicant tracking systems
- 📱 Responsive web interface built with Angular
- 🔒 Secure API built with Go and Gin framework
- 🗃️ Persistent data storage with PostgreSQL and pgvector
- 🔄 Real-time updates and instant preview
- 📤 Export to PDF format

## 🚀 Getting Started

### Prerequisites

- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)
- An OpenRouter API key for AI capabilities

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

## 🔗 Accessing the Application

- Frontend: http://localhost:4200
- Backend API: http://localhost:8080/api
- Health Check: http://localhost:8080/health

## 📚 API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/resume` | Get all resumes |
| GET | `/api/resume/:id` | Get a resume by ID |
| POST | `/api/resume` | Create a new resume |
| PUT | `/api/resume/:id` | Update a resume |
| DELETE | `/api/resume/:id` | Delete a resume |
| GET | `/api/skills` | Get all skills |
| GET | `/api/experience` | Get all experiences |
| POST | `/api/chat/message` | Send a message to the chatbot |
| GET | `/api/chat/history/:sessionId` | Get chat history |
| POST | `/api/chat/generate-resume` | Generate a resume from chat history |

## 📁 Project Structure

```
resume.in/
├── backend/              # Go backend
│   ├── config/           # Configuration
│   ├── controllers/      # HTTP handlers
│   ├── middleware/       # Middleware
│   ├── models/           # Data models
│   ├── routes/           # API routes
│   ├── utils/            # Utilities
│   └── test/             # Test files and output
├── client/               # Angular frontend
│   ├── src/              # Source code
│   └── ...
├── docker-compose.yml    # Docker Compose configuration
├── start.sh              # Script to start the application (Linux/Mac)
├── stop.sh               # Script to stop the application (Linux/Mac)
├── start.bat             # Script to start the application (Windows)
└── stop.bat              # Script to stop the application (Windows)
```

## 💻 Development

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

## 🤝 Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

Please make sure to update tests as appropriate and adhere to the existing coding style.

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Gin Framework](https://github.com/gin-gonic/gin) for the Go web framework
- [Angular](https://angular.io/) for the frontend framework
- [pgvector](https://github.com/pgvector/pgvector) for vector similarity search in PostgreSQL
- [OpenRouter](https://openrouter.ai/) for providing AI capabilities

## 📧 Contact

If you have any questions or suggestions, please open an issue or reach out to the maintainers.
