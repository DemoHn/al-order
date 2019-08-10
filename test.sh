#!/bin/sh
ENV_FILE=.env
set -a
[ -f $ENV_FILE ] && source $ENV_FILE

# run test commands
DATABASE_URL=$DATABASE_URL go test ./...

set +a