#!/usr/bin/env bash

if [ -z $1 ]; then
    echo "pass in migration name"
    exit 1
fi

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
go run github.com/golang-migrate/migrate/v4/cmd/migrate create -ext sql -dir "$DIR/../migrations" $1