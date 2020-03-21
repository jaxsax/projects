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

    cd "$BUILD_WORKSPACE_DIRECTORY"
    bazelisk build //tapeworm/botv2/cmd/bot:docker.tar
    docker load -i bazel-bin/tapeworm/botv2/cmd/bot/docker.tar
    generatedImageID=$(docker images  bazel/tapeworm/botv2/cmd/bot:docker --format '{{.ID}}')
    version="v$(date '+%Y%m%d-%H%M')"
    docker tag "$generatedImageID" r.internal.jaxsax.co/tapeworm/botv2:$version
    docker push r.internal.jaxsax.co/tapeworm/botv2:$version
)
