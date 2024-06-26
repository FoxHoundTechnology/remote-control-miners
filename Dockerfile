# Start from the official Golang base image version 1.20 on Alpine Linux
FROM golang:1.20

# Environment variables which CompileDaemon requires to run
ENV PROJECT_DIR=/app \
    GO111MODULE=on \
    CGO_ENABLED=0

# Basic setup of the container
RUN mkdir /app
COPY .. /app

# Establish /app as the working directory within the container
WORKDIR /app

# Get CompileDaemon
RUN go get github.com/githubnemo/CompileDaemon
RUN go install github.com/githubnemo/CompileDaemon

# Transfer the Go module dependency files into the container
COPY go.mod go.sum ./

# Fetch and install the dependencies defined in the Go module files, including the missing module
RUN go mod download && \
    go mod tidy

# Transfer the application source code into the container
COPY . .

# Compile the Go application and output the executable as 'main'
RUN go build -o main .

# Make port 8080 available to the world outside this container
EXPOSE 8080

# Set the command that will be executed when the container starts
# CMD ["./main"]
ENTRYPOINT CompileDaemon -build="go build -o remote-control-server" -command="./remote-control-server"
