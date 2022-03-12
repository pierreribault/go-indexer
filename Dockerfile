FROM golang:alpine as Builder
WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download
COPY cmd cmd
COPY pkg pkg

RUN go build -o /go-indexer /app/cmd/main.go

ENTRYPOINT ["/go-indexer"]

FROM alpine:3.15
WORKDIR /app

COPY --from=Builder /go-indexer /go-indexer

ENTRYPOINT ["/go-indexer"]