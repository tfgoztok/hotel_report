FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /hotel-service ./cmd/api

FROM alpine:3.14

WORKDIR /

COPY --from=builder /hotel-service /hotel-service
COPY --from=builder /app/internal/db/migrations /internal/db/migrations

# Copy test files
COPY --from=builder /app/tests /tests

EXPOSE 8080

CMD ["/hotel-service"]