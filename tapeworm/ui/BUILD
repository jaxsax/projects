load("@build_bazel_rules_nodejs//:index.bzl", "nodejs_binary")
load("@npm_bazel_jasmine//:index.bzl", "jasmine_node_test")
load("@npm//http-server:index.bzl", "http_server")
load("@build_bazel_rules_nodejs//internal/web_package:web_package.bzl", "web_package")

load(":parcel.bzl", "parcel")

parcel(
    name = "bundle",
    srcs = [
        "src/bar.js"
    ],
    entry_point = "src/foo.js"
)

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
    deps = [
        ":bundle",
    ]
)

web_package(
    name = "stuff",
    index_html = "index.html",
    data = [
        ":bundle"
    ]
)

http_server(
    name = "statics",
    data = [
        ":stuff",
    ],
    templated_args = ["tapeworm/ui/stuff"],
)