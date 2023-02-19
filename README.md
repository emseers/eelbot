# Eelbot

[![Build Status](https://github.com/emseers/eelbot/actions/workflows/go.yml/badge.svg)](https://github.com/emseers/eelbot/actions)

Eelbot is a simple bot that can listen to commands on do things. Eelbot needs a
PostgreSQL database to store the various required data to function. Notable
features include:

* Posting random (or specific) jokes (from the database)
* Posting random (or specific) pictures (from the database)
* Posting random (or specific) audio taunts (from the database)

## Setup

Eelbot requires PostgreSQL 15.2 or later (earlier versions may work, but is
untested). The schema definition is in the [initdb](initdb/) directory. If using
[Docker](https://www.docker.com/), you should be able to use it as a volume
mount for the [initialization of the schema](https://github.com/docker-library/docs/tree/master/postgres#initialization-scripts).
Once initialized, use [eelbot-ui](https://github.com/emseers/eelbot-ui) to
prep/modify the database. Modifying the database directly is not recommended.

## Build

Install [golang](https://golang.org/) 1.20 or later. Clone the project and build
with:

```
mkdir bin
go build -o ./bin/ ./cmd/...
```

## Run

Grab the config template from `configs/eelbot/config.yaml` and update it to
match your setup, then save it. Run the bot executable in the `bin` folder:

```
./eelbot -c <path/to/config/file> -t <discord-bot-token>
```

Note that passing the path to the config is optional, and if it's not passed it
looks for a file called either `config.json` or `config.yaml` in the current
folder.

## Docker

As an alternative to building and running eelbot locally, you can build and run
it in [Docker](https://www.docker.com/). Build the image with the provided
Dockerfile:

```
docker build -t eelbot:latest .
```

You need to mount all required files and folders to run the container:

```
docker run \
    -it \
    --name eelbot \
    -v <full/path/to/config.yaml>:/app/config.yaml \
    eelbot:latest \
    -t <discord-bot-token>
```
