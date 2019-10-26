load("@build_bazel_rules_nodejs//:index.bzl", "nodejs_binary")

nodejs_binary(
    name = "server",
    entry_point = "run.js",
    data = [
        "@npm//express",
        "@npm//moment",
    ],
)
