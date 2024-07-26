# Start from the official Golang base image version 1.21
FROM golang:1.21

# Environment variables
ENV GO111MODULE=on \
    CGO_ENABLED=0

# Set the working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Install a specific version of Air compatible with Go 1.21
RUN go install github.com/cosmtrek/air@v1.49.0

# Copy the Air configuration file
COPY .air.toml ./

# Copy the rest of the application
COPY . .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run Air (which will build and run your app, watching for changes)
CMD ["air"]