load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "user",
    srcs = [
        "user.go",
        "where.go",
    ],
    importpath = "github.com/m-nny/universe/ent/user",
    visibility = ["//visibility:public"],
    deps = [
        "//ent/predicate",
        "@io_entgo_ent//dialect/sql",
        "@io_entgo_ent//dialect/sql/sqlgraph",
    ],
)
