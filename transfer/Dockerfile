FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go install github.com/swaggo/swag/cmd/swag@latest && swag init -g cmd/main.go
RUN go build -o transport ./cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/transport .
COPY --from=builder /app/docs ./docs

EXPOSE 8080
CMD ["./transport"]