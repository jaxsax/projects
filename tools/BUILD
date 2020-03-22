package(default_visibility = ["//visibility:public"])

sh_binary(
    name = "buildifier",
    srcs = ["bazel_tools/buildifier.sh"],
    data = [
        "@buildifier//file",
    ],
)

sh_binary(
    name = "cache-gc",
    srcs = ["sys/gc-cache.sh"],
)

sh_binary(
    name = "test",
    srcs = ["test.sh"],
)

sh_binary(
    name = "update-bot",
    srcs = ["update-bot.sh"],
)

sh_binary(
    name = "update-deps",
    srcs = ["update-deps.sh"],
)
