load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

def dependencies():
    http_archive(
        name = "io_bazel_rules_docker",
        sha256 = "1698624e878b0607052ae6131aa216d45ebb63871ec497f26c67455b34119c80",
        strip_prefix = "rules_docker-0.15.0",
        urls = ["https://github.com/bazelbuild/rules_docker/releases/download/v0.15.0/rules_docker-v0.15.0.tar.gz"],
    )

    # Gazelle
    http_archive(
        name = "io_bazel_rules_go",
        urls = [
            "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.24.13/rules_go-v0.24.13.tar.gz",
            "https://github.com/bazelbuild/rules_go/releases/download/v0.24.13/rules_go-v0.24.13.tar.gz",
        ],
        sha256 = "52d0a57ea12139d727883c2fef03597970b89f2cc2a05722c42d1d7d41ec065b",
    )

    http_archive(
        name = "bazel_gazelle",
        urls = [
            "https://mirror.bazel.build/github.com/bazelbuild/bazel-gazelle/releases/download/v0.22.3/bazel-gazelle-v0.22.3.tar.gz",
            "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.22.3/bazel-gazelle-v0.22.3.tar.gz",
        ],
        sha256 = "222e49f034ca7a1d1231422cdb67066b885819885c356673cb1f72f748a3c9d4",
    )

    http_archive(
        name = "build_bazel_rules_nodejs",
        sha256 = "3653eb344b7222189f781e260522acf77195f5d7ef643bc23ab3a738445c1394",
        urls = ["https://github.com/bazelbuild/rules_nodejs/releases/download/4.5.0/rules_nodejs-4.5.0.tar.gz"],
    )
