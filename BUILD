exports_files(
    ["tsconfig.json"],
    visibility = ["//visibility:public"],
)

load("@bazel_gazelle//:def.bzl", "gazelle")
load("@py_deps//:requirements.bzl", "requirement")

# gazelle:prefix github.com/jaxsax/projects
gazelle(name = "gazelle")

py_library(
    name = 'pytest',
    deps = [
        requirement('py'),
        requirement('packaging'),
        requirement('pluggy'),
        requirement('importlib-metadata'),
        requirement('zipp'),
        requirement('attrs'),
        requirement('more-itertools'),
        requirement('wcwidth'),
        requirement('pytest'),
    ],
    visibility = ["//visibility:public"],
)
