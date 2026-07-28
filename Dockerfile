# Build Stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy dependency definition
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Final Stage (Minimal Alpine Image)
FROM alpine:latest  

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy compiled binary & migration files from builder stage
COPY --from=builder /app/main .
COPY --from=builder /app/migrations ./migrations

# Default Port
EXPOSE 8080

CMD ["./main"]
