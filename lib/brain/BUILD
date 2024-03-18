load("@rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "brain",
    srcs = [
        "album.go",
        "artist.go",
        "brain.go",
        "track.go",
    ],
    importpath = "github.com/m-nny/universe/lib/brain",
    visibility = ["//visibility:public"],
    deps = [
        "//lib/utils/sliceutils",
        "@com_github_zmb3_spotify_v2//:spotify",
        "@io_gorm_driver_sqlite//:sqlite",
        "@io_gorm_gorm//:gorm",
    ],
)

go_test(
    name = "brain_test",
    srcs = [
        "album_test.go",
        "artist_test.go",
        "brain_test.go",
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
