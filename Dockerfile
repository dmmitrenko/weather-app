FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o weatherapp ./cmd/main.go

FROM alpine:latest
RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=builder /app/weatherapp .
COPY --from=builder /app/configs/config.yaml ./configs/config.yaml

EXPOSE 8080
CMD ["./weatherapp"]
