load("@io_bazel_rules_docker//container:container.bzl", "container_pull")

def repositories():
    container_pull(
        name = "ubuntu-1910-base",
        registry = "index.docker.io",
        repository = "ubuntu",
        tag = "19.10",
        digest = "sha256:7ce552ad1c3e94a5c3d2bb24c07000c34a4bb43fd9b379652b2c80593a018e80",
    )

    container_pull(
        name = "ubuntu-2004-base",
        registry = "index.docker.io",
        repository = "ubuntu",
        tag = "20.04",
        digest = "sha256:3093096ee188f8ff4531949b8f6115af4747ec1c58858c091c8cb4579c39cc4e",
    )

    container_pull(
        name = "ubuntu-2010-base",
        registry = "index.docker.io",
        repository = "ubuntu",
        tag = "20.10",
        digest = "sha256:160a9181d622d428f6836e17245fea90b87e9f7abb86939d002c2e301383c8a8",
    )

    container_pull(
        name = "base-images-ubuntu-2004",
        registry = "r.internal.jaxsax.co",
        repository = "base-images/ubuntu",
        tag = "20.04",
        digest = "sha256:a0a79e1af3fb8de1fd41847b21bed4c026e47bf6776ce4c4ca499b11c6bef9cc",
    )

