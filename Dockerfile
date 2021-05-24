FROM golang:1.13-alpine AS build-env

RUN apk add --no-cache git mercurial build-base

RUN mkdir -p /go/src/github.com/emseers/eelbot
WORKDIR /go/src/github.com/emseers/eelbot
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .
RUN mkdir bin && go build -o ./bin/ ./cmd/...

FROM alpine:3.11

RUN apk update && apk add --no-cache ca-certificates

COPY --from=build-env /go/src/github.com/emseers/eelbot/bin/eelbot /app/eelbot

WORKDIR /app
ENTRYPOINT ["./eelbot"]
