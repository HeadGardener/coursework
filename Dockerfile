FROM golang:1.21.0

WORKDIR /app

COPY ./ ./

RUN go mod download
RUN go build -o coursework ./cmd/coursework/main.go

EXPOSE 5000

CMD ["./coursework"]
