load("@rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "brain",
    srcs = [
        "artist.go",
        "brain.go",
        "brainIndex.go",
        "common.go",
        "discogs_release.go",
        "meta_album.go",
        "meta_track.go",
        "spotify_album.go",
        "spotify_track.go",
        "user.go",
    ],
    importpath = "github.com/m-nny/universe/lib/brain",
    visibility = ["//visibility:public"],
    deps = [
        "//lib/discogs",
        "//lib/utils/sliceutils",
        "//lib/utils/spotifyutils",
        "@com_github_zmb3_spotify_v2//:spotify",
        "@io_gorm_driver_sqlite//:sqlite",
        "@io_gorm_gorm//:gorm",
        "@io_gorm_gorm//logger",
        "@org_golang_x_oauth2//:oauth2",
    ],
)

go_test(
    name = "brain_test",
    srcs = [
        "artist_test.go",
        "brain_test.go",
        "meta_album_test.go",
        "meta_track_test.go",
    ],
    embed = [":brain"],
    deps = [
        "@com_github_google_go_cmp//cmp",
        "@com_github_google_go_cmp//cmp/cmpopts",
        "@com_github_zmb3_spotify_v2//:spotify",
    ],
)
