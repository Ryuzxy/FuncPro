FROM golang:1.25-alpine

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o main ./cmd

# Expose port
EXPOSE 8080

# Run the application
CMD ["./main"]