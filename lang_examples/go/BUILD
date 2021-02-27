load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library", "go_test")
load("@io_bazel_rules_docker//go:image.bzl", "go_image")
load("@io_bazel_rules_docker//container:container.bzl", "container_image", "container_push")

go_image(
    name = "image_base",
    srcs = ["main.go"],
)

container_image(
    name = "image",
    base = ":image_base",
    ports = ["9090"],
)

go_test(
    name = "go_test",
    srcs = ["main_test.go"],
    embed = [":go_lib"],
)

go_binary(
    name = "main",
    embed = [":go_lib"],
    visibility = ["//visibility:public"],
)

go_library(
    name = "go_lib",
    srcs = ["main.go"],
    importpath = "github.com/jaxsax/projects/lang_examples/go",
    visibility = ["//visibility:private"],
)
