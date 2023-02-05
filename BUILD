# gazelle:prefix github.com/smocherla-brex/proto-lint-unused-imports
load("@bazel_gazelle//:def.bzl", "gazelle")
load("@com_github_bazelbuild_buildtools//buildifier:def.bzl", "buildifier")

gazelle(name = "gazelle")

gazelle(
    name = "gazelle-update-repos",
    args = [
        "-from_file=go.mod",
        "-to_macro=go_deps.bzl%go_deps",
        "-prune",
    ],
    command = "update-repos",
)

buildifier(
    name = "buildifier",
)
