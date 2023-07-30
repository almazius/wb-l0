FROM golang:latest

WORKDIR /wb/

COPY .. .

RUN go mod download
EXPOSE 3000

RUN go build -o main cmd/main.go

CMD ["./main"]