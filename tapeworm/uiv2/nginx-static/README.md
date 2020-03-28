Illustrates how to build a static image docker container using nginx to serve static files. The main use case of this will be have an external nginx reverse proxying to this container for its contents whilst serving an API through another path.

# Usage Instructions

```
bazelisk build :static_files_image.tar
docker load -i  ../../../bazel-bin/tapeworm/uiv2/nginx-static/static_files_image.tar
docker run -p 9999:80 --name nginxstatic -d --rm  bazel/tapeworm/uiv2/nginx-static:static_files_image
```
