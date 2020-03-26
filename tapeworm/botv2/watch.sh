#!/usr/bin/env bash

if [[ -n "${BUILD_WORKSPACE_DIRECTORY:-}" ]]; then
    echo "You're running via bazel!" >&2
elif ! command -v bazelisk &>/dev/null; then
    echo "Install bazelisk at https://github.com/bazelbuild/bazelisk" >&2
    exit 1
fi

WEB_PORT=9999 ibazel run //tapeworm/botv2/cmd/bot:bot -- -config_path=$PWD/secrets/config.yml 
