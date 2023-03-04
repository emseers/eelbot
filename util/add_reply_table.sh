#!/usr/bin/env bash
set -e

host="${POSTGRESQL_HOST:-localhost}"
port="${POSTGRESQL_PORT:-5432}"
user="${POSTGRESQL_USER}"
password="${POSTGRESQL_PASSWORD}"
database="${POSTGRESQL_DATABASE}"

if [ -z "${user}" ]; then
	conn_str="postgresql://${host}:${port}"
elif [ -z "${password}" ]; then
	conn_str="postgresql://${user}@${host}:${port}"
else
	conn_str="postgresql://${user}:${password}@${host}:${port}"
fi
if ! [ -z "${database}" ]; then
	conn_str+="/${database}"
fi

read -p "Name for reply type: " name

command=$(cat <<- END
	CREATE TABLE reply.$name (
		id   integer PRIMARY KEY,
		text text NOT NULL
	);
END
)

psql -c "$command" "$conn_str"
