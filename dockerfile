# ====================
# Dockerfile (optional)
# ====================
FROM golang:1.21-alpine

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go install github.com/fsnotify/fsnotify@latest

CMD ["go", "run", "thunder.go"]
