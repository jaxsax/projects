package(default_visibility = ["//visibility:public"])

sh_binary(
    name = "buildifier",
    srcs = ["bazel_tools/buildifier.sh"],
    data = [
        "@buildifier//file",
    ]
)
