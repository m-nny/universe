load("@rules_go//go:def.bzl", "go_library")

go_library(
    name = "spotify",
    srcs = [
        "albums_clean.go",
        "artist_clean.go",
        "spotify.go",
        "tracks_clean.go",
        "user_tracks_clean.go",
    ],
    importpath = "github.com/m-nny/universe/lib/spotify",
    visibility = ["//visibility:public"],
    deps = [
        "//lib/brain",
        "//lib/jsoncache",
        "//lib/spotify/token",
        "//lib/utils/sliceutils",
        "@com_github_zmb3_spotify_v2//:spotify",
        "@com_github_zmb3_spotify_v2//auth",
    ],
)
