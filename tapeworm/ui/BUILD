load("@build_bazel_rules_nodejs//:index.bzl", "nodejs_binary")
load("@npm_bazel_jasmine//:index.bzl", "jasmine_node_test")

nodejs_binary(
    name = "server",
    data = [
        "@npm//express",
        "@npm//moment",
    ],
    entry_point = "run.js",
)

jasmine_node_test(
    name = "test",
    srcs = glob(["*.spec.js"]),
)
