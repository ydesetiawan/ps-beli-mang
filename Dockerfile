# Use the official Golang base image with version 1.21.4
FROM golang:1.21.4-alpine as builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go application files into the container
COPY . .

# Build the Go application
ENV GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target=/root/.cache/go-build GOARCH=amd64 GOOS=linux go build -o main ./cmd/api

# Use a minimal base image for the final image
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the built executable from the builder stage
COPY --from=builder /app/main .

# Expose the port the application will run on
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
