# Build stage
FROM golang:1.23-alpine AS builder


RUN apk add --no-cache git ca-certificates
WORKDIR /app

ENV GOPROXY=https://proxy.golang.org,direct

COPY go.mod go.sum ./
RUN cat go.mod && echo "----" && go env && echo "----" && go mod download -x

COPY . .
RUN go build -o payout-api

# Run stage
FROM alpine:latest
RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=builder /app/payout-api .

EXPOSE 8080
CMD ["./payout-api"]
