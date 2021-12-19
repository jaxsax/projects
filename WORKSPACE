# gazelle:repository_macro repos.bzl%go_repositories
workspace(
    name = "com_github_jaxsax_projects",
    managed_directories = {"@npm": ["node_modules"]},
)

load("//:dependencies.bzl", "dependencies")

dependencies()

load(
    "@bazel_tools//tools/build_defs/repo:http.bzl",
    "http_archive",
    "http_file",
)
load("@bazel_tools//tools/build_defs/repo:git.bzl", "git_repository")

git_repository(
    name = "com_google_protobuf",
    commit = "ae50d9b9902526efd6c7a1907d09739f959c6297",
    remote = "https://github.com/protocolbuffers/protobuf",
    shallow_since = "1613677815 -0800",
)

load("@com_google_protobuf//:protobuf_deps.bzl", "protobuf_deps")

protobuf_deps()

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains(version = "1.16.6")

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

load("@build_bazel_rules_nodejs//:index.bzl", "node_repositories", "yarn_install")

node_repositories(
    node_version = "16.13.0",
    package_json = ["//:package.json"],
)

yarn_install(
    # Name this npm so that Bazel Label references look like @npm//package
    name = "npm",
    package_json = "//:package.json",
    quiet = False,
    yarn_lock = "//:yarn.lock",
)

# Tools
http_file(
    name = "buildifier",
    executable = True,
    sha256 = "4c985c883eafdde9c0e8cf3c8595b8bfdf32e77571c369bf8ddae83b042028d6",
    urls = ["https://github.com/bazelbuild/buildtools/releases/download/0.29.0/buildifier"],
)

load("@rules_pkg//:deps.bzl", "rules_pkg_dependencies")

rules_pkg_dependencies()

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")

gazelle_dependencies(go_repository_default_config = "//:WORKSPACE.bazel")
