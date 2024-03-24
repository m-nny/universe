load("@rules_go//go:def.bzl", "go_library")

go_library(
    name = "spotify_ent",
    srcs = [
        "albums_ent.go",
        "artist_ent.go",
        "spotify.go",
        "token_ent.go",
        "tracks_ent.go",
    ],
    importpath = "github.com/m-nny/universe/lib/spotify_ent",
    visibility = ["//visibility:public"],
    deps = [
        "//ent",
        "//ent/album",
        "//ent/artist",
        "//ent/track",
        "//ent/user",
        "//lib/jsoncache",
        "//lib/utils/hitcounter",
        "//lib/utils/maputils",
        "//lib/utils/sliceutils",
        "//lib/utils/spotifyutils",
        "@com_github_zmb3_spotify_v2//:spotify",
        "@com_github_zmb3_spotify_v2//auth",
        "@org_golang_x_oauth2//:oauth2",
    ],
)
