load("@rules_go//go:def.bzl", "go_library")

go_library(
    name = "brain",
    srcs = [
        "artist.go",
        "brain.go",
    ],
    importpath = "github.com/m-nny/universe/lib/brain",
    visibility = ["//visibility:public"],
    deps = [
        "//lib/utils/sliceutils",
        "@com_github_zmb3_spotify_v2//:spotify",
        "@io_gorm_gorm//:gorm",
    ],
)
