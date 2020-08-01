# gazelle:exclude wire.go
load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "main.go",
        "wire_gen.go",
    ],
    importpath = "github.com/jaxsax/projects/lang_examples/go/wiredemo",
    visibility = ["//visibility:private"],
)

go_binary(
    name = "wire",
    embed = [":go_default_library"],
    visibility = ["//visibility:public"],
)
