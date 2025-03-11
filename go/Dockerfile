# Use golang:1.24-alpine as the base image for building the application
FROM golang:1.24-alpine AS builder

# Enable CGO for go-sqlite3
ENV CGO_ENABLED=1

# Install necessary build dependencies for CGO and SQLite
RUN apk add --no-cache gcc musl-dev sqlite

# Set the working directory inside the container to `/app`
WORKDIR /app

# Create a user 'trainee' and group 'mercari' to run the application securely
RUN addgroup -S mercari && adduser -S trainee -G mercari

# Copy only go.mod and go.sum first to take advantage of Docker layer caching
COPY go.mod go.sum ./

# Download all dependencies based on go.mod and go.sum
RUN go mod download

# Copy the remaining application source code into the container
COPY . .

# Change ownership of the db directory (now at /app/db)
RUN chown -R trainee:mercari db

# Grant trainee user read & write permissions for db directory
RUN chmod -R u+rw db

# Build the application binary and output it as '/app/mercari-app'
RUN go build -o /app/mercari-app ./cmd/api

USER trainee

# Set the command to execute the application when the container starts
CMD ["/app/mercari-app"]
