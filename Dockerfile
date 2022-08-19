FROM golang:1.19.0

ENV GO111MODULE=on \
    GOOS=linux \
    GOARCH=amd64

ADD main.go /go/src/main.go
WORKDIR /go/src/

RUN apt update && apt install -y sqlite3 gcc

RUN mkdir -p data

RUN go mod init gobank
RUN go mod tidy

RUN go build -o main .

CMD ["./main"]