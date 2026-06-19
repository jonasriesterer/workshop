FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

FROM alpine:3

WORKDIR /app

COPY --from=builder /app/server ./server
COPY --from=builder /app/web    ./web
COPY --from=builder /app/db     ./db

EXPOSE 3000

CMD ["./server"]
