FROM golang:alpine AS build-env

RUN apk add --no-cache git mercurial build-base

COPY . /go/src/github.com/Emseers/Eelbot
WORKDIR /go/src/github.com/Emseers/Eelbot
RUN go build

FROM alpine:3.11

RUN apk update && apk add --no-cache ca-certificates

COPY --from=build-env /go/src/github.com/Emseers/Eelbot/Eelbot /executable/Eelbot
COPY EelbotDB.db /executable/EelbotDB.db
COPY config.ini /executable/config.ini
COPY taunts /executable/taunts
COPY pics /executable/pics

WORKDIR /executable
ENTRYPOINT ./Eelbot -t $EELBOT_TOKEN