FROM golang:1.24.4-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./

RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/api/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app/

# Create a non-root user
RUN addgroup -g 1000 kinfam && \
    adduser -D -u 1000 -G kinfam kinfam

# Switch to non-root user
USER kinfam
COPY --from=builder /app/main /app/
COPY --from=builder /app/api/docs/ /app/api/docs/

EXPOSE 8080
CMD ["./main"]
