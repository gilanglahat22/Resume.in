@echo off
echo Starting Resume.in application...

REM Check if Docker is installed
where docker >nul 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo Docker is not installed. Please install Docker first.
    exit /b 1
)

REM Check if Docker Compose is installed
where docker-compose >nul 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo Docker Compose is not installed. Please install Docker Compose first.
    exit /b 1
)

REM Build and start the containers
echo Building and starting Docker containers...
docker-compose up --build

REM Wait for user input to exit
echo Press Ctrl+C to stop the application. 