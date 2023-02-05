load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

go_library(
    name = "unused_imports_lib",
    srcs = ["unused_imports.go"],
    importpath = "github.com/smocherla-brex/proto-lint-unused-imports/internal/unused_imports",
    visibility = ["//visibility:private"],
    deps = [
        "@com_github_bufbuild_protocompile//:protocompile",
        "@com_github_bufbuild_protocompile//ast",
        "@com_github_bufbuild_protocompile//reporter",
    ],
)

go_binary(
    name = "unused_imports",
    embed = [":unused_imports_lib"],
    visibility = ["//:__subpackages__"],
)
