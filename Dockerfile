# Stage 1: Build
FROM golang:alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -o main .

# Stage 2: Final image
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy the binary from the builder
COPY --from=builder /app/main .
# Copy docs if needed (swagger files are usually embedded or in a folder)
COPY --from=builder /app/docs ./docs

# Expose port
EXPOSE 8090

# Start the application
CMD ["./main"]
