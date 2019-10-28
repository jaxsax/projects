load("@npm_bazel_typescript//:index.bzl", "ts_library", "ts_devserver")


ts_library(
    name = "widget",
    srcs = glob(["widget/*.ts"])
)

ts_library(
    name = "code",
    srcs = [
        "index.ts",
    ],
    deps = [
        ":widget",
    ]
)

ts_devserver(
    name = "devserver",
    port = 8080,
    deps = [
        ":code"
    ],
    static_files = [
        "index.html",
        "src/main.css",
    ],
    serving_path = "/bundle.js"
)
