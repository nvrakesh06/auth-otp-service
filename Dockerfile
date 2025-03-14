# Use Go base image
FROM golang:1.23

# Set working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum first (for dependency caching)
COPY go.mod go.sum ./

# Download dependencies
RUN go mod tidy

# Copy the rest of the application source code
COPY . .

# Build the application
RUN go build -o otp-service ./cmd/main.go

# Expose the application port
EXPOSE 8080

# Run the application
CMD [ "./otp-service" ]
