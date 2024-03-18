load("@rules_go//go:def.bzl", "go_library")

go_library(
    name = "spotify",
    srcs = [
        "albums.go",
        "artist.go",
        "spotify.go",
        "token.go",
        "tracks.go",
        "user_tracks.go",
        "utils.go",
    ],
    importpath = "github.com/m-nny/universe/lib/spotify",
    visibility = ["//visibility:public"],
    deps = [
        "//ent",
        "//ent/album",
        "//ent/artist",
        "//ent/track",
        "//ent/user",
        "//lib/brain",
        "//lib/jsoncache",
        "//lib/utils/hitcounter",
        "//lib/utils/maputils",
        "//lib/utils/sliceutils",
        "@com_github_zmb3_spotify_v2//:spotify",
        "@com_github_zmb3_spotify_v2//auth",
        "@org_golang_x_oauth2//:oauth2",
    ],
)
