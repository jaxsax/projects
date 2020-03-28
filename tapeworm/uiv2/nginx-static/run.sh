#!/usr/bin/env bash

set -xe

bazelisk build :static_files_image.tar
docker load -i  ../../../bazel-bin/tapeworm/uiv2/nginx-static/static_files_image.tar
docker run -p 9999:80 --name nginxstatic -d --rm  bazel/tapeworm/uiv2/nginx-static:static_files_image
