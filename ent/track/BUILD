load("@rules_go//go:def.bzl", "go_library")

go_library(
    name = "track",
    srcs = [
        "track.go",
        "where.go",
        "where_custom.go",
    ],
    importpath = "github.com/m-nny/universe/ent/track",
    visibility = ["//visibility:public"],
    deps = [
        "//ent/predicate",
        "@io_entgo_ent//dialect/sql",
        "@io_entgo_ent//dialect/sql/sqlgraph",
        "@io_entgo_ent//dialect/sql/sqljson",
    ],
)
