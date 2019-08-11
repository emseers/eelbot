FROM golang:alpine

RUN apk add --no-cache git mercurial build-base \
&& go get github.com/bwmarrin/discordgo \
&& go get github.com/mattn/go-sqlite3

WORKDIR /go/src