workspace(name = "proto-lint-unused-imports")

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

http_archive(
    name = "io_bazel_rules_go",
    sha256 = "56d8c5a5c91e1af73eca71a6fab2ced959b67c86d12ba37feedb0a2dfea441a6",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.37.0/rules_go-v0.37.0.zip",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.37.0/rules_go-v0.37.0.zip",
    ],
)

http_archive(
    name = "bazel_gazelle",
    sha256 = "ecba0f04f96b4960a5b250c8e8eeec42281035970aa8852dda73098274d14a1d",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-gazelle/releases/download/v0.29.0/bazel-gazelle-v0.29.0.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.29.0/bazel-gazelle-v0.29.0.tar.gz",
    ],
)

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")
load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")

############################################################
# Define your own dependencies here using go_repository.
# Else, dependencies declared by rules_go/gazelle will be used.
# The first declaration of an external repository "wins".
############################################################

load("//:go_deps.bzl", "go_deps")

# gazelle:repository_macro go_deps.bzl%go_deps
go_deps()

go_rules_dependencies()

go_register_toolchains(version = "1.19.5")

gazelle_dependencies()

http_archive(
    name = "com_google_protobuf",
    sha256 = "023e2bb164b234af644c5049c6dac1d9c9f6dd2acb133b960d9009105b4226bd",
    strip_prefix = "protobuf-27.4",
    # latest, as of 2020-02-21
    urls = [
        "https://mirror.bazel.build/github.com/protocolbuffers/protobuf/archive/v27.4.tar.gz",
        "https://github.com/protocolbuffers/protobuf/archive/v27.4.tar.gz",
    ],
)

load("@com_google_protobuf//:protobuf_deps.bzl", "protobuf_deps")

protobuf_deps()
