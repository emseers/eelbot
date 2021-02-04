# Eelbot

Eelbot is a simple bot that can listen to commands on do things. Eelbot needs an SQLite database to store the various required data to function. Notable features include:

* Posting random jokes (from the database)
* Posting random pictures (whose paths are stored in the database)
* Posting random audio taunts (from a designated folder)

## Setup

Eelbot requires a SQLite database with the following schema:
```sql
CREATE TABLE sqlite_sequence(name,seq);
CREATE TABLE `EelPics` (
	`ImageID`	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	`CheckSum`	TEXT NOT NULL,
	`FullPath`	TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS "EelJokes" (
	`JokeID`	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	`JokeText`	TEXT NOT NULL,
	`JokeTextLine2`	TEXT
);
CREATE TABLE `CmdInvokes` (
	`ID`	INTEGER NOT NULL,
	`Invoker`	TEXT NOT NULL,
	`TimeStamp`	REAL NOT NULL,
	`CmdID`	INTEGER NOT NULL,
	PRIMARY KEY(`ID`)
);
CREATE TABLE `Cmds` (
	`ID`	INTEGER NOT NULL,
	`Title`	TEXT NOT NULL UNIQUE,
	`Description`	TEXT,
	PRIMARY KEY(`ID`)
);
```

### Jokes

Once you have the SQLite database setup, to use the jokes functionality (with `/badjoke me`), add jokes to the `EelJokes` table. Single line jokes only require the `JokeText` column filled out and multiline jokes require the punchline to be in stored in `JokeTextLine2` as well.

### Pics

To use the pictures functionality (with `/eel me` or `/eel bomb <X>`), add pictures to the `EelPics` table. At this time, only the `FullPath` column is required (`CheckSum` can be left blank).

### Taunts

To use the taunts functionality (with `/taunt <X>`), create a `taunts` folder and store all the taunt files in there. Note that this folder should only contain audio files. The taunt number is inferred from the files in the order they are in when sorted alphabetically, so you may want to rename the files to get the numbers you wish.

## Build

Install [golang](https://golang.org/) 1.13 or later. Clone the project and build with:

```
mkdir bin
go build -o ./bin/ ./cmd/...
```

## Run

Grab the config template from `configs/eelbot/config.ini` and update it to match your setup, then save it. Run the bot executable in the `bin` folder:

```
./eelbot -c <path/to/config/file> -t <discord-bot-token>
```

Note that passing the path to the config is optional, and if it's not passed it looks for a file called `config.ini` in the current folder.

## Docker

As an alternative to building and running Eelbot locally, you can build and run it in [Docker](https://www.docker.com/). Build the image with the provided Dockerfile:

```
docker build -t eelbot:latest .
```

You need to mount all required files and folders to run the container:

```
docker run \
    -it \
    --name eelbot \
    -v <full/path/to/EelbotDB.db>:/app/EelbotDB.db \
    -v <full/path/to/config.ini>:/app/config.ini \
    -v <full/path/to/pics/folder>:/app/pics \
    -v <full/path/to/taunts/folder>:/app/taunts \
    eelbot:latest \
    -t <discord-bot-token>
```
