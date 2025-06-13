# Build stage
FROM golang:1.20-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o lc2lt .

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/lc2lt .

# Create .aws directory
RUN mkdir -p /root/.aws

# Set environment variables
ENV AWS_PROFILE=qa
ENV AWS_REGION=us-west-2

# Run the application
CMD ["./lc2lt"] 