# =====================================
# Build Stage - ใช้สำหรับ compile Go application
# =====================================
FROM golang:1.23-alpine AS builder

# ติดตั้ง dependencies ที่จำเป็นสำหรับ build
RUN apk add --no-cache git ca-certificates

# สร้าง working directory
WORKDIR /app

# Copy go mod และ go sum files เพื่อ download dependencies
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code ทั้งหมด
COPY . .

# Build the Go app สำหรับ Linux
# CGO_ENABLED=0 เพื่อสร้าง static binary
# GOOS=linux เพื่อให้แน่ใจว่า build สำหรับ Linux
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# =====================================
# Runtime Stage - ใช้สำหรับ run application
# =====================================
FROM alpine:latest AS runtime

# ติดตั้ง ca-certificates สำหรับ HTTPS requests
RUN apk --no-cache add ca-certificates

# สร้าง user ที่ไม่ใช่ root เพื่อความปลอดภัย
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# สร้าง working directory
WORKDIR /root/

# Copy binary จาก build stage
COPY --from=builder /app/main .

# Copy ca-certificates จาก build stage ถ้าจำเป็น
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# เปลี่ยน ownership ของไฟล์ให้กับ appuser
RUN chown -R appuser:appgroup /root/

# เปลี่ยนไปใช้ non-root user
USER appuser

# Expose port (ปรับตามที่ใช้ใน application)
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/health || exit 1

# Run the binary
CMD ["./main"]