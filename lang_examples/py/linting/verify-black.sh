#!/usr/bin/env bash

export LC_ALL=C.UTF-8
export LANG=C.UTF-8

set -o errexit
set -o nounset
set -o pipefail

DIR=$( cd "$( dirname "$0" )" && pwd )

if [[ -n "${TEST_WORKSPACE:-}" ]]; then # Running inside bazel
  pwd
  echo "$DIR"
  echo "Linting python..." >&2
elif ! command -v bazel &> /dev/null; then
  echo "Install bazel at https://bazel.build" >&2
  exit 1
else
  (
    set -o xtrace
    bazel test --test_output=streamed //experiments:verify-black
  )
  exit 0
fi

shopt -s extglob globstar

"$DIR/black_bin" --check !(external)/**/*.py
