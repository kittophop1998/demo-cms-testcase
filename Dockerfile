# ---------------------------------------------------------
# Stage 1: Build Go binary
# ---------------------------------------------------------
FROM golang:1.25-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go.mod/go.sum and download dependencies
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy the entire project
COPY . .

# Move to working dir for main.go
WORKDIR /app/cmd/server

# Build Go binary
RUN go build -o server .

# ---------------------------------------------------------
# Stage 2: Run with Google Chrome
# ---------------------------------------------------------
FROM debian:stable-slim

# Install certificates (for HTTPS support)
# RUN apk --no-cache add ca-certificates

# Set working directory
WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/cmd/server/server .

# (Optional) Copy configs if needed
COPY configs/ /root/configs/

# Expose port (เปลี่ยนตามที่ Go app ใช้)
EXPOSE 8080

# Run the binary
CMD ["./server"]
