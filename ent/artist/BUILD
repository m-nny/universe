load("@rules_go//go:def.bzl", "go_library")

go_library(
    name = "artist",
    srcs = [
        "artist.go",
        "where.go",
    ],
    importpath = "github.com/m-nny/universe/ent/artist",
    visibility = ["//visibility:public"],
    deps = [
        "//ent/predicate",
        "@io_entgo_ent//dialect/sql",
        "@io_entgo_ent//dialect/sql/sqlgraph",
    ],
)
