#!/bin/bash

if [ -z "$1" ]; then
	echo "Eelbot token not supplied"
	exit 1
fi

if [ "$2" == "-b" ]; then
	docker build -t eelbot-build:v1 build
	docker run --name eelbot-build --rm -itd -v $(readlink -f .):/go/src/eelbot eelbot-build:v1
	docker exec eelbot-build sh -c 'cd eelbot && go build'
	docker stop eelbot-build
fi

docker build -t eelbot:v1 .
docker run --name eelbot --rm -e EELBOT_TOKEN=$1 -it eelbot:v1