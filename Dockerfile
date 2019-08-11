FROM alpine:3.10
RUN apk update && apk add --no-cache ca-certificates

COPY eelbot /executable/eelbot
COPY EelbotDB.db /executable/EelbotDB.db

WORKDIR /executable
ENTRYPOINT ./eelbot -t $EELBOT_TOKEN