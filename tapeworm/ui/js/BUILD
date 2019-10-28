load("//tapeworm/ui:parcel.bzl", "parcel")
load("@build_bazel_rules_nodejs//:index.bzl", "nodejs_binary")
load("@npm_bazel_jasmine//:index.bzl", "jasmine_node_test")
load("@build_bazel_rules_nodejs//internal/web_package:web_package.bzl", "web_package")
load("@npm//http-server:index.bzl", "http_server")

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
    assets = [
        "src/main.css"
    ],
    data = [
        ":bundle"
    ]
)

http_server(
    name = "js",
    data = [
        ":stuff",
    ],
    templated_args = ["tapeworm/ui/stuff"],
)
