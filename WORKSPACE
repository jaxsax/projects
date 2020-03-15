workspace(
    name = "com_github_jaxsax_projects",
    managed_directories = {"@npm": ["node_modules"]},
)

load(
    "@bazel_tools//tools/build_defs/repo:http.bzl",
    "http_archive",
    "http_file",
)

http_archive(
    name = "rules_python",
    sha256 = "aa96a691d3a8177f3215b14b0edc9641787abaaa30363a080165d06ab65e1161",
    url = "https://github.com/bazelbuild/rules_python/releases/download/0.0.1/rules_python-0.0.1.tar.gz",
)

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
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.22.1/rules_go-v0.22.1.tar.gz",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.22.1/rules_go-v0.22.1.tar.gz",
    ],
    sha256 = "e6a6c016b0663e06fa5fccf1cd8152eab8aa8180c583ec20c872f4f9953a7ac5",
)

http_archive(
    name = "bazel_gazelle",
    urls = [
        "https://storage.googleapis.com/bazel-mirror/github.com/bazelbuild/bazel-gazelle/releases/download/v0.20.0/bazel-gazelle-v0.20.0.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.20.0/bazel-gazelle-v0.20.0.tar.gz",
    ],
    sha256 = "d8c45ee70ec39a57e7a05e5027c32b1576cc7f16d9dd37135b0eddde45cf1b10",
)

load("@io_bazel_rules_go//go:deps.bzl", "go_rules_dependencies", "go_register_toolchains")

go_rules_dependencies()

go_register_toolchains()

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies", "go_repository")

gazelle_dependencies()

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

local_repository(
    name = "thirdparty",
    path = "third_party",
)

# Node
http_archive(
    name = "build_bazel_rules_nodejs",
    sha256 = "26c39450ce2d825abee5583a43733863098ed29d3cbaebf084ebaca59a21a1c8",
    urls = ["https://github.com/bazelbuild/rules_nodejs/releases/download/0.39.0/rules_nodejs-0.39.0.tar.gz"],
)

load("@build_bazel_rules_nodejs//:index.bzl", "node_repositories", "yarn_install")

node_repositories(
    node_version = "10.12.0",
    vendored_node = "@thirdparty//:node-v13.0.1-linux-x64",
    vendored_yarn = "@thirdparty//:yarn-v1.19.1",
)

yarn_install(
    name = "npm",
    data = [
        "@thirdparty//:node-v13.0.1-linux-x64/bin/node",
        "@thirdparty//:yarn-v1.19.1/bin/yarn.js",
    ],
    package_json = "//:package.json",
    yarn_lock = "//:yarn.lock",
)

load("@npm//:install_bazel_dependencies.bzl", "install_bazel_dependencies")

install_bazel_dependencies()

# Set up TypeScript toolchain
load("@npm_bazel_typescript//:index.bzl", "ts_setup_workspace")

ts_setup_workspace()

# Python
load("@rules_python//python:repositories.bzl", "py_repositories")

py_repositories()

# Only needed if using the packaging rules.
load("@rules_python//python:pip.bzl", "pip_import", "pip_repositories")

pip_repositories()

pip_import(
    name = "tapeworm_import",
    requirements = "//tapeworm/bot:requirements.txt",
)

load(
    "@tapeworm_import//:requirements.bzl",
    _tapeworm_install = "pip_install",
)

_tapeworm_install()

# Tools
http_file(
    name = "buildifier",
    executable = True,
    sha256 = "4c985c883eafdde9c0e8cf3c8595b8bfdf32e77571c369bf8ddae83b042028d6",
    urls = ["https://github.com/bazelbuild/buildtools/releases/download/0.29.0/buildifier"],
)

go_repository(
    name = "com_github_go_telegram_bot_api_telegram_bot_api",
    importpath = "github.com/go-telegram-bot-api/telegram-bot-api",
    sum = "h1:2cauKuaELYAEARXRkq2LrJ0yDDv1rW7+wrTEdVL3uaU=",
    version = "v4.6.4+incompatible",
)

go_repository(
    name = "com_github_technoweenie_multipartstreamer",
    importpath = "github.com/technoweenie/multipartstreamer",
    sum = "h1:XRztA5MXiR1TIRHxH2uNxXxaIkKQDeX7m2XsSOlQEnM=",
    version = "v1.0.1",
)
