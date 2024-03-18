load("@rules_go//go:def.bzl", "go_library")

go_library(
    name = "utils",
    srcs = [
        "album.go",
        "artist.go",
        "track.go",
    ],
    importpath = "github.com/m-nny/universe/lib/spotify/utils",
    visibility = ["//visibility:public"],
    deps = [
        "//lib/utils/sliceutils",
        "@com_github_zmb3_spotify_v2//:spotify",
    ],
)
