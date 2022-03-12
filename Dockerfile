FROM golang:alpine
WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download
COPY cmd cmd
COPY pkg pkg

RUN go build -o /go-indexer /app/cmd/main.go

ENTRYPOINT ["/go-indexer"]