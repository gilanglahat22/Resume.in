@echo off
echo Stopping any running containers...
docker-compose down

echo Removing PostgreSQL volume to get a fresh start...
FOR /F "tokens=*" %%G IN ('docker volume ls -q ^| findstr postgres-data') DO (
    docker volume rm %%G
)

echo Starting containers...
docker-compose up -d

echo Containers started, you can check logs with: docker-compose logs -f 