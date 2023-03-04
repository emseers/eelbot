# Eelbot

[![Build Status](https://github.com/emseers/eelbot/actions/workflows/go.yml/badge.svg)](https://github.com/emseers/eelbot/actions)

Eelbot is a simple bot that can listen to commands on do things. Eelbot needs a
PostgreSQL database to store the various required data to function. Notable
features include:

* Posting random (or specific) jokes (from the database)
* Posting random (or specific) media files (from the database)
* Posting replies (from the database) on regex matches

## Setup

Eelbot requires PostgreSQL 15.2 or later (earlier versions may work, but is
untested). The schema definition is in the [initdb](initdb/) directory. If using
[Docker](https://www.docker.com/), you should be able to use it as a volume
mount for the [initialization of the schema](https://github.com/docker-library/docs/tree/master/postgres#initialization-scripts).
Once initialized, use [eelbot-ui](https://github.com/emseers/eelbot-ui) to
prep/modify the database. Alternatively, steps for manually updating the
database is described below.

### Jokes

The `jokes` table in the `public` schema contains simple jokes. Jokes can either
be a single line, or a joke with a leadup and a punchline. For single line
jokes, only the `text` needs to be set. Otherwise, both `text` and `punchline`
can be set to make it a joke with a punchline. The `id` can be any unique value,
but is required to contain gapless sequential values. This table can safely be
modified at runtime without requiring a restart.

Jokes can be posted using the `/badjoke me` command for a random joke or
`/badjoke <id>` for a specific joke.

### Images/Taunts

The `images` and `taunts` tables in the `public` schema have the same format.
Their difference is in name only (and by extension their trigger command) and
otherwise work exactly the same. An entry can be added by adding the bytes of a
`file` with a corresponding `name`. The `id` can be any unique value, but is
required to contain gapless sequential values. These tables can safely be
modified at runtime without requiring a restart.

Images can be posted using the `/eel me` command for a random file or
`/eel <id>` for a specific file. Taunts can be posted using the `/taunt me`
command for a random file or `/taunt <id>` for a specific file.

### Replies

Custom replies can be triggered on regex matches, but requires some setup. For
each type of reply, there needs to be:

* one or more regular expressions specified that determine if it's a match
* one or more values to choose from for the reply
* a trigger chance (percentage) on whether a reply will be sent on a match
* a min and max treshold for time to wait before sending the reply
* a cooldown period during which this reply can't be triggered again

The following steps are required to add a reply type:

1. Create a table to contain the set of reply values in the `reply` schema:

   ```psql
   CREATE TABLE reply.$name (id integer PRIMARY KEY, text text NOT NULL);
   ```

   where `$name` is the name of the reply type.
1. Add the set of reply values to the table. The `text` can be any string and
   the `id` can be any unique value, but is required to contain gapless
   sequential values.
1. Add an entry to the config file under the `replies` section with the
   following values:

   * `table_name`: the name of the reply type (should match the name of the
     table created in the `reply` schema)
   * `regexps`: a list of regular expressions to use to determine if the reply
     should be triggered (values should be valid golang regex flavor)
   * `percent`: the trigger chance for the reply (0 - 100)
   * `min_delay`: minimum time to wait before sending the reply
   * `max_delay`: maximum time to wait before sending the reply
   * `timeout`: timeout (in seconds) for consequetive replies of this type

The reply type should now be able to be triggered by receiving any text that
matches one of the match regex values. The values in the config file are loaded
once at startup, but the values in the custom tables in the `reply`
schema are read lazily at runtime. Therefore, any changes to the config file
will require a restart to take effect, whereas the contents of the tables in the
`reply` schema can safely be modified at runtime without requiring a restart.

## Build

Install [golang](https://golang.org/) 1.20 or later. Clone the project and build
with:

```sh
mkdir bin
go build -o ./bin/ ./cmd/...
```

## Run

Grab the config template from `configs/eelbot/config.yaml` and update it to
match your setup, then save it. Run the bot executable in the `bin` folder:

```sh
./eelbot -c <path/to/config/file> -t <discord-bot-token>
```

Note that passing the path to the config is optional, and if it's not passed it
looks for a file called either `config.json` or `config.yaml` in the current
folder.

## Docker

As an alternative to building and running eelbot locally, you can build and run
it in [Docker](https://www.docker.com/). Build the image with the provided
Dockerfile:

```sh
docker build -t eelbot:latest .
```

You need to mount all required files and folders to run the container:

```sh
docker run \
    -it \
    --name eelbot \
    -v <full/path/to/config.yaml>:/app/config.yaml \
    eelbot:latest \
    -t <discord-bot-token>
```
