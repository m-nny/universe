load("@rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "brain",
    srcs = [
        "album.go",
        "artist.go",
        "brain.go",
        "brain.sqlx.go",
        "brainIndex.go",
        "common.go",
        "discogs_release.go",
        "discogs_seller.go",
        "meta_album.go",
        "meta_track.go",
        "spotify_album.go",
        "spotify_track.go",
        "track.go",
        "user.go",
    ],
    importpath = "github.com/m-nny/universe/lib/brain",
    visibility = ["//visibility:public"],
    deps = [
        "//lib/discogs",
        "//lib/utils/sliceutils",
        "//lib/utils/spotifyutils",
        "@com_github_jmoiron_sqlx//:sqlx",
        "@com_github_tursodatabase_libsql_client_go//libsql",
        "@com_github_zmb3_spotify_v2//:spotify",
        "@org_golang_x_oauth2//:oauth2",
    ],
)

go_test(
    name = "brain_test",
    srcs = [
        "album_test.go",
        "artist_test.go",
        "brain_test.go",
        "meta_album_test.go",
        "meta_track_test.go",
        "spotify_album_test.go",
        "spotify_track_test.go",
        "track_test.go",
    ],
    embed = [":brain"],
    deps = [
        "@com_github_google_go_cmp//cmp",
        "@com_github_google_go_cmp//cmp/cmpopts",
        "@com_github_jmoiron_sqlx//:sqlx",
        "@com_github_mattn_go_sqlite3//:go-sqlite3",
        "@com_github_zmb3_spotify_v2//:spotify",
    ],
)
