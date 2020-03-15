load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["main.go"],
    importpath = "github.com/jaxsax/projects/tapeworm/bot_v2",
    visibility = ["//visibility:private"],
    deps = ["@com_github_go_telegram_bot_api_telegram_bot_api//:go_default_library"],
)

go_binary(
    name = "bot_v2",
    embed = [":go_default_library"],
    visibility = ["//visibility:public"],
)
