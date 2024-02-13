load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "schema",
    srcs = [
        "album.go",
        "artist.go",
        "track.go",
        "user.go",
    ],
    importpath = "github.com/m-nny/universe/ent/schema",
    visibility = ["//visibility:public"],
    deps = [
        "@io_entgo_ent//:ent",
        "@io_entgo_ent//schema/edge",
        "@io_entgo_ent//schema/field",
        "@org_golang_x_oauth2//:oauth2",
    ],
)
