# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o payout-api

# Run stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/payout-api .

EXPOSE 8080
CMD ["./backend"]
