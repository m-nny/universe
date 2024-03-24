load("@rules_go//go:def.bzl", "go_library")

go_library(
    name = "spotify_ent",
    srcs = [
        "albums_ent.go",
        "artist_ent.go",
        "spotify.go",
        "tracks_ent.go",
        "user_ent.go",
    ],
    importpath = "github.com/m-nny/universe/lib/spotify_ent",
    visibility = ["//visibility:public"],
    deps = [
        "//ent",
        "//ent/album",
        "//ent/artist",
        "//ent/track",
        "//lib/utils/hitcounter",
        "//lib/utils/maputils",
        "//lib/utils/sliceutils",
        "//lib/utils/spotifyutils",
        "@com_github_zmb3_spotify_v2//:spotify",
    ],
)
