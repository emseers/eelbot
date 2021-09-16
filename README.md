# Eelbot

Eelbot is a simple bot that can listen to commands on do things. Eelbot needs an
SQLite database to store the various required data to function. Notable features
include:

* Posting random (or specific) jokes (from the database)
* Posting random (or specific) pictures (whose paths are stored in the database)
* Posting random (or specific) audio taunts (whose paths are stored in the
  database)

## Setup

Eelbot requires a SQLite database with the following schema:
```sql
CREATE TABLE sqlite_sequence(name,seq);
CREATE TABLE "images" (
  "id"   INTEGER NOT NULL,
  "path" TEXT NOT NULL,
  PRIMARY KEY("id" AUTOINCREMENT)
);
CREATE TABLE "jokes" (
  "id"        INTEGER NOT NULL,
  "text"      TEXT NOT NULL,
  "punchline" TEXT,
  PRIMARY KEY("id" AUTOINCREMENT)
);
CREATE TABLE "taunts" (
  "id"   INTEGER NOT NULL,
  "path" TEXT NOT NULL,
  PRIMARY KEY("id" AUTOINCREMENT)
);
```

### Jokes

To use the jokes functionality (with `/badjoke me` or `/badjoke <X>`), add jokes
to the `jokes` table. Single line jokes only require the `text` column and
multiline jokes require the punchline to be in stored in `punchline` as well.

### Images

To use the images functionality (with `/eel me` or `/eel <X>`), add pictures to
the `images` table. The `path` should contain the full path to the image.

### Taunts

To use the taunts functionality (with `/taunt me` or `/taunt <X>`), add taunts
to the `taunts` table. The `path` should contain the full path to the taunt.

## Build

Install [golang](https://golang.org/) 1.17 or later. Clone the project and build
with:

```
mkdir bin
go build -o ./bin/ ./cmd/...
```

## Run

Grab the config template from `configs/eelbot/config.ini` and update it to match
your setup, then save it. Run the bot executable in the `bin` folder:

```
./eelbot -c <path/to/config/file> -t <discord-bot-token>
```

Note that passing the path to the config is optional, and if it's not passed it
looks for a file called `config.ini` in the current folder.

## Docker

As an alternative to building and running Eelbot locally, you can build and run
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
    -v <full/path/to/eelbot.db>:/app/eelbot.db \
    -v <full/path/to/config.ini>:/app/config.ini \
    -v <full/path/to/images/folder>:/app/images \
    -v <full/path/to/taunts/folder>:/app/taunts \
    eelbot:latest \
    -t <discord-bot-token>
```
