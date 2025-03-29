# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/main ./main.go

# Final image
FROM alpine:3.19

WORKDIR /app
RUN apk --no-cache add ca-certificates && \
    addgroup -g 1000 appuser && \
    adduser -u 1000 -G appuser -D appuser

COPY --from=builder --chown=appuser:appuser /app/main /app/main

USER appuser
EXPOSE 8001
CMD ["/app/main"]