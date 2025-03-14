# Use the official Go image
FROM golang:1.23

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files and install dependencies
COPY go.mod go.sum ./
RUN go mod tidy

# Copy the rest of the application code
COPY . .

# Install Fiber & other dependencies
RUN go build -o otp-service

# Expose the application port
EXPOSE 8080

# Command to run the application
CMD ["./otp-service"]
