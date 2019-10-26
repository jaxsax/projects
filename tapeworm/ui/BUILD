load("@build_bazel_rules_nodejs//:index.bzl", "nodejs_binary")

nodejs_binary(
    name = "server",
    data = [
        "@npm//express",
        "@npm//moment",
    ],
    entry_point = "run.js",
)
