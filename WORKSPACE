# gazelle:repository_macro repos.bzl%go_repositories
workspace(
    name = "com_github_jaxsax_projects",
)

load(
    "@bazel_tools//tools/build_defs/repo:http.bzl",
    "http_archive",
    "http_file",
)
load("@bazel_tools//tools/build_defs/repo:git.bzl", "git_repository")

# git_repository(
#     name = "rules_python",
#     remote = "https://github.com/bazelbuild/rules_python.git",
#     commit = "748aa53d7701e71101dfd15d800e100f6ff8e5d1",
#     shallow_since = "1583438240 -0500",
# )

# http_archive(
#     name = "build_bazel_rules_nodejs",
#     sha256 = "591d2945b09ecc89fde53e56dd54cfac93322df3bc9d4747cb897ce67ba8cdbf",
#     urls = ["https://github.com/bazelbuild/rules_nodejs/releases/download/1.2.0/rules_nodejs-1.2.0.tar.gz"],
# )

# http_archive(
#     name = "build_bazel_rules_svelte",
#     url = "https://github.com/thelgevold/rules_svelte/archive/0.15.zip",
#     strip_prefix = "rules_svelte-0.15",
#     sha256 = "1b04eb08ef80636929d152bb2f2733e36d9e0b8ad10aca7b435c82bd638336f5",
# )

http_archive(
    name = "io_bazel_rules_docker",
    sha256 = "14ac30773fdb393ddec90e158c9ec7ebb3f8a4fd533ec2abbfd8789ad81a284b",
    strip_prefix = "rules_docker-0.12.1",
    urls = ["https://github.com/bazelbuild/rules_docker/releases/download/v0.12.1/rules_docker-v0.12.1.tar.gz"],
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

# Python
#load("@rules_python//python:repositories.bzl", "py_repositories")
#
#py_repositories()
#
## Only needed if using the packaging rules.
#load("@rules_python//python:pip.bzl", "pip3_import")
#
#pip3_import(
#    name = "py_deps",
#    requirements = "//:requirements.txt",
#)
#
#load("@py_deps//:requirements.bzl", "pip_install")
#
#pip_install()

load("@io_bazel_rules_go//go:deps.bzl", "go_rules_dependencies", "go_register_toolchains")

go_rules_dependencies()

go_register_toolchains()

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")

gazelle_dependencies()

git_repository(
    name = "com_google_protobuf",
    commit = "ae50d9b9902526efd6c7a1907d09739f959c6297",
    remote = "https://github.com/protocolbuffers/protobuf",
    shallow_since = "1613677815 -0800"
)

load("@com_google_protobuf//:protobuf_deps.bzl", "protobuf_deps")

protobuf_deps()

load("//:repos.bzl", "go_repositories")

go_repositories()

load("//:containers.bzl", "repositories")
repositories()

load(
    "@io_bazel_rules_docker//repositories:repositories.bzl",
    container_repositories = "repositories",
)

container_repositories()

load(
    "@io_bazel_rules_docker//go:image.bzl",
    _go_image_repos = "repositories",
)

_go_image_repos()



# Javascript

# load("@build_bazel_rules_nodejs//:index.bzl", "yarn_install")

# yarn_install(
#     name = "npm",
#     package_json = "//:package.json",
#     yarn_lock = "//:yarn.lock",
# )

# load("@npm//:install_bazel_dependencies.bzl", "install_bazel_dependencies")

# install_bazel_dependencies()

# load("@build_bazel_rules_svelte//:defs.bzl", "rules_svelte_dependencies")

# rules_svelte_dependencies()

# Tools
http_file(
    name = "buildifier",
    executable = True,
    sha256 = "4c985c883eafdde9c0e8cf3c8595b8bfdf32e77571c369bf8ddae83b042028d6",
    urls = ["https://github.com/bazelbuild/buildtools/releases/download/0.29.0/buildifier"],
)

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")
http_archive(
    name = "rules_pkg",
    url = "https://github.com/bazelbuild/rules_pkg/releases/download/0.2.5/rules_pkg-0.2.5.tar.gz",
    sha256 = "352c090cc3d3f9a6b4e676cf42a6047c16824959b438895a76c2989c6d7c246a",
)
load("@rules_pkg//:deps.bzl", "rules_pkg_dependencies")
rules_pkg_dependencies()
