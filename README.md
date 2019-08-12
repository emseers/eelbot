# Eelbot

## Build

Install [golang](https://golang.org/) 1.12 or later. Get dependencies with:

    go get github.com/bwmarrin/discordgo
    go get github.com/mattn/go-sqlite3

Clone project to `$GOPATH/src/eelbot`. Build with:

    go build

## Run

Run the bot with the token:

    ./eelbot -t <EELBOT_TOKEN>

## Docker

Build and run the bot in an alpine image with the provided Dockerfile:

    docker build -t eelbot:v1 .
    docker run --name eelbot --rm -e EELBOT_TOKEN=<EELBOT_TOKEN> -it eelbot:v1