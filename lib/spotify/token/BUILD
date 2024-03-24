load("@rules_go//go:def.bzl", "go_library")

go_library(
    name = "token",
    srcs = [
        "brain_token_storage.go",
        "inmemory_token_storage.go",
        "token_storage.go",
    ],
    importpath = "github.com/m-nny/universe/lib/spotify/token",
    visibility = ["//visibility:public"],
    deps = [
        "//lib/brain",
        "@com_github_zmb3_spotify_v2//auth",
        "@org_golang_x_oauth2//:oauth2",
    ],
)
