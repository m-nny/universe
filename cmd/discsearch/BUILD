load("@rules_go//go:def.bzl", "go_binary", "go_library")

go_library(
    name = "discsearch_lib",
    srcs = ["main.go"],
    importpath = "github.com/m-nny/universe/cmd/discsearch",
    visibility = ["//visibility:private"],
    deps = [
        "//ent/album",
        "//ent/artist",
        "//lib/discsearch",
        "//lib/spotify",
        "//lib/spotify/utils",
        "@io_entgo_ent//dialect/sql",
    ],
)

go_binary(
    name = "discsearch",
    embed = [":discsearch_lib"],
    visibility = ["//visibility:public"],
)
