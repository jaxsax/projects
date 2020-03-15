#!/usr/bin/env bash

if [[ -n "${BUILD_WORKSPACE_DIRECTORY:-}" ]]; then
    echo "You're running via bazel!" >&2
elif ! command -v bazelisk &>/dev/null; then
    echo "Install bazelisk at https://github.com/bazelbuild/bazelisk" >&2
    exit 1
else
    (
        set -o xtrace
        bazelisk run //tools:test
    )
    exit 0
fi

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

echo "pwd: $PWD"
echo "dir: $DIR"
