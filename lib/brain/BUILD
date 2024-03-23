load("@rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "brain",
    srcs = [
        "albumIndex.go",
        "artist.go",
        "brain.go",
        "common.go",
        "meta_album.go",
        "spotify_album.go",
        "track.go",
    ],
    importpath = "github.com/m-nny/universe/lib/brain",
    visibility = ["//visibility:public"],
    deps = [
        "//lib/spotify/utils",
        "//lib/utils/sliceutils",
        "@com_github_zmb3_spotify_v2//:spotify",
        "@io_gorm_driver_sqlite//:sqlite",
        "@io_gorm_gorm//:gorm",
        "@io_gorm_gorm//logger",
    ],
)

go_test(
    name = "brain_test",
    srcs = [
        "artist_test.go",
        "brain_test.go",
        "spotify_album_test.go",
        "track_test.go",
    ],
    embed = [":brain"],
    deps = [
        "@com_github_google_go_cmp//cmp",
        "@com_github_google_go_cmp//cmp/cmpopts",
        "@com_github_zmb3_spotify_v2//:spotify",
        "@io_gorm_gorm//:gorm",
    ],
)
