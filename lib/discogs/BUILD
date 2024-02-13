load("@rules_go//go:def.bzl", "go_library")

go_library(
    name = "discogs",
    srcs = [
        "client.go",
        "marketplace.go",
        "master.go",
        "models.go",
        "release.go",
    ],
    importpath = "github.com/m-nny/universe/lib/discogs",
    visibility = ["//visibility:public"],
    deps = ["//lib/jsoncache"],
)
