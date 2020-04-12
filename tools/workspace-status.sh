#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

git_commit="$(git describe --tags --always --dirty)"
build_date="$(date -u '+%Y%m%d')"
docker_tag="v${build_date}-${git_commit}"
cat <<EOF
DOCKER_REPO ${DOCKER_REPO:-localhost:5000}
DOCKER_TAG ${docker_tag}
EOF