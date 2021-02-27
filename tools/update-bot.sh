#!/usr/bin/env bash

if [[ -n "${BUILD_WORKSPACE_DIRECTORY:-}" ]]; then
    echo "You're running via bazel!" >&2
elif ! command -v bazelisk &>/dev/null; then
    echo "Install bazelisk at https://github.com/bazelbuild/bazelisk" >&2
    exit 1
else
    (
        set -o xtrace
        bazelisk run //tools:update-bot
    )
    exit 0
fi

if ! command -v docker &>/dev/null; then
    echo "Requires docker"
    exit 1
fi

(
    set -o xtrace
    set -e

    cd "$BUILD_WORKSPACE_DIRECTORY"
    # bazelisk test //tapeworm/botv2/...
    bazelisk run //tapeworm/botv2/cmd/bot:push
)
