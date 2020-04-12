load("@io_bazel_rules_docker//container:container.bzl", "container_pull")

def repositories():
    container_pull(
        name = "ubuntu-1910-base",
        registry = "index.docker.io",
        repository = "ubuntu",
        tag = "19.10",
        digest = "sha256:7ce552ad1c3e94a5c3d2bb24c07000c34a4bb43fd9b379652b2c80593a018e80",
    )
