FROM golang:1.18 AS setup
WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download
COPY cmd cmd
COPY pkg pkg

RUN go build -o /go-indexer /app/cmd/watch/main.go

CMD ["/go-indexer"]

# Build Go-Indexer
#
FROM setup AS builder
ARG CGO_ENABLED=0
RUN apt-get update \
    && apt-get install -y upx \
    && go build -ldflags "-s -w" -o /go-indexer /app/cmd/watch/main.go \
    && upx --best --lzma /go-indexer

# Store Go-Indexer in a scratch image
#
FROM alpine:3.15 AS production
COPY --from=builder /go-indexer /go-indexer
CMD [ "/go-indexer" ]