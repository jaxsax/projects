load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library", "go_test")
load("@io_bazel_rules_docker//go:image.bzl", "go_image")
load("@io_bazel_rules_docker//container:container.bzl", "container_push", "container_image")

go_image(
    name = "image_base",
    srcs = ["main.go"]
)

container_image(
    name = "image",
    base = ":image_base",
    ports = ["9090"]
)

go_test(
    name = "go_default_test",
    srcs = ["main_test.go"]
)

go_binary(
  name = "main",
  srcs = ["main.go"],
)
