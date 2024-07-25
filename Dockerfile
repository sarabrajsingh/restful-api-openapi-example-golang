# Use the official Golang image as the base image
FROM golang:1.21.6-alpine AS build

# Set the working directory inside the container
WORKDIR /app

# Install git for go mod download
RUN apk add --no-cache git

# Copy the Go modules manifests
COPY go.mod go.sum ./

# Download the Go modules
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application
RUN go build -o app-api

# Use a minimal Alpine image for the final container
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Install necessary CA certificates
RUN apk add --no-cache ca-certificates

# Copy the Go binary from the build stage
COPY --from=build /app/app-api .

# Copy the config and other necessary files
COPY --from=build /app/api /app/api
COPY --from=build /app/swaggerui/dist /app/swaggerui/dist

# Expose the port the application runs on
EXPOSE 8080

# Command to run the Go binary
CMD ["/app/app-api"]
