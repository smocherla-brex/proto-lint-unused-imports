load("@rules_proto//proto:defs.bzl", "proto_library")
load("@io_bazel_rules_go//go:def.bzl", "go_library")
load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library")
load("//:lint.bzl", "proto_unused_imports_test")

proto_library(
    name = "simple_proto",
    srcs = [
        "bar.proto",
        "foo.proto",
    ],
    visibility = ["//visibility:public"],
    deps = [
        "@com_google_protobuf//:any_proto",
        "@com_google_protobuf//:descriptor_proto",
    ],
)

go_proto_library(
    name = "simple_go_proto",
    importpath = "github.com/smocherla-brex/proto-lint-unused-imports/examples/simple",
    proto = ":simple_proto",
    visibility = ["//visibility:public"],
)

go_library(
    name = "simple",
    embed = [":simple_go_proto"],
    importpath = "github.com/smocherla-brex/proto-lint-unused-imports/examples/simple",
    visibility = ["//visibility:public"],
)

proto_unused_imports_test(
    name = "simple_proto_test",
    protos = [":simple_proto"],
)
