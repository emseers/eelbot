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

Run the bot in an alpine image using the run script:

    ./run.sh <EELBOT_TOKEN> -b