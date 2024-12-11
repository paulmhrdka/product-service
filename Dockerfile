FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install required packages
RUN apk add --no-cache gcc musl-dev

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application and migration tool
RUN CGO_ENABLED=0 GOOS=linux go build -o api ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux go build -o migrate ./cmd/migration

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy binaries from builder
COPY --from=builder /app/api .
COPY --from=builder /app/migrate .
COPY --from=builder /app/data ./data

# Copy wait-for script to ensure Elasticsearch is ready
COPY scripts/wait-for.sh .
RUN chmod +x wait-for.sh

EXPOSE 8080
