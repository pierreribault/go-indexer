FROM golang:alpine
WORKDIR /app

ADD go.mod .
ADD go.sum .

RUN go mod download
ADD cmd cmd
ADD pkg pkg

RUN go build -o /go-indexer /app/cmd/main.go

ENTRYPOINT ["/go-indexer"]