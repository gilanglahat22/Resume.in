FROM golang:1.18 AS builder

WORKDIR /app

# Copy go mod file first
COPY go.mod ./

# Initialize go modules
RUN go mod download

# Copy the rest of the source code
COPY . .

# Run go mod tidy to ensure dependencies are correct
RUN go mod tidy

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Use a smaller image for the final application
FROM alpine:latest

WORKDIR /app

# Install necessary packages
RUN apk --no-cache add ca-certificates tzdata

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Expose the port the app runs on
EXPOSE 8080

# Command to run the application
CMD ["./main"] 