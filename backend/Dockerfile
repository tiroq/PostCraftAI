# Build stage
FROM golang:1.20-alpine AS builder
WORKDIR /app
COPY go.mod ./
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 go build -o service .

# Run stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/service .
EXPOSE 8080
CMD ["./service"]
