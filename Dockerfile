FROM golang:1.22.0-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY ./ ./

RUN go build -o coursework ./cmd/coursework/main.go

FROM alpine

WORKDIR /app

COPY --from=builder /app/coursework /app/coursework

EXPOSE 8080

CMD ["./coursework"]
