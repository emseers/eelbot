FROM golang:alpine AS build-env

RUN apk add --no-cache git mercurial build-base \
&& go get github.com/bwmarrin/discordgo \
&& go get github.com/mattn/go-sqlite3

COPY . /go/src/eelbot/
WORKDIR /go/src/eelbot
RUN go build

FROM alpine:3.10

COPY --from=build-env /go/src/eelbot/eelbot /executable/eelbot
COPY EelbotDB.db /executable/EelbotDB.db

WORKDIR /executable
ENTRYPOINT ./eelbot -t $EELBOT_TOKEN