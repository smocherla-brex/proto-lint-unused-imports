load("@rules_proto//proto:defs.bzl", "proto_library")
load("@io_bazel_rules_go//go:def.bzl", "go_library")
load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library")
load("//:lint.bzl", "proto_unused_imports_test")

proto_library(
    name = "examples_complex_package2_proto",
    srcs = ["bar.proto"],
    visibility = ["//visibility:public"],
    deps = ["//examples/complex/package1:examples_complex_package1_proto"],
)

go_proto_library(
    name = "examples_complex_package2_go_proto",
    importpath = "github.com/smocherla-brex/proto-lint-unused-imports/examples/complex/package2",
    proto = ":examples_complex_package2_proto",
    visibility = ["//visibility:public"],
    deps = ["//examples/complex/package1"],
)

go_library(
    name = "package2",
    embed = [":examples_complex_package2_go_proto"],
    importpath = "github.com/smocherla-brex/proto-lint-unused-imports/examples/complex/package2",
    visibility = ["//visibility:public"],
)

proto_unused_imports_test(
    name = "examples_complex_package2_proto_test",
    protos = [":examples_complex_package2_proto"],
)