# Stage 1: Build the Go application
FROM golang:1.23-alpine AS build
WORKDIR /app
RUN go install github.com/air-verse/air@latest
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
COPY . .
RUN go mod download
RUN sqlc generate -f /app/sql/sqlc.yaml
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o server .

# set executable
CMD ["air", "-c", ".air.toml"]