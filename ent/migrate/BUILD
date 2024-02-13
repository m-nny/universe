load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "migrate",
    srcs = [
        "migrate.go",
        "schema.go",
    ],
    importpath = "github.com/m-nny/universe/ent/migrate",
    visibility = ["//visibility:public"],
    deps = [
        "@io_entgo_ent//dialect",
        "@io_entgo_ent//dialect/sql/schema",
        "@io_entgo_ent//schema/field",
    ],
)
