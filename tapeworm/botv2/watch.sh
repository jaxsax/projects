#!/usr/bin/env bash

err=0

$(command -v bazelisk &>/dev/null) || {
	echo "Please install bazelisk at https://github.com/bazelbuild/bazelisk" >&2
	err=1
}

$(command -v ibazel &>/dev/null) || {
	echo "Please install ibazel at https://github.com/bazelbuild/bazel-watcher" >&2
	err=1
}

$(command -v docker-compose &> /dev/null) || {
	echo "Please install docker-compose at https://github.com/docker/compose" >&2
	err=1
}

if [ $err -eq 1 ]; then
	exit 1
fi

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd)"
mkdir -p $DIR/secrets/
secrets_config=$DIR/secrets/config.yml

if [ ! -f "$secrets_config" ]; then
	echo "Creating default config, please populate telegram bot token"
	cat <<-EOF > "$secrets_config"
		token: "default_token"
		database:
            sqlite_db_path: database.db
	EOF
fi

WEB_PORT=9999 ibazel run //tapeworm/botv2/cmd/bot:bot -- -config_path="$secrets_config"
