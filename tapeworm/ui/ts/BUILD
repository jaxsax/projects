load("@npm_bazel_typescript//:index.bzl", "ts_library", "ts_devserver")
load("@npm_bazel_rollup//:index.bzl", "rollup_bundle")


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
    ],
    serving_path = "/bundle.js"
)

ts_library(
    name = "app",
    srcs = [
        "app.ts"
    ],
    deps = [
        ":widget"
    ]
)

rollup_bundle(
    name = "app_chunks",
    deps = [
        ":app"
    ],
    sourcemap = "inline",
    entry_point = "index",
    output_dir = True
)

ts_devserver(
    name = "app_devserver",
    deps = [
        ":app"
    ]
)