load("@rules_go//go:def.bzl", "go_library")

go_library(
    name = "discsearch",
    srcs = [
        "discsearch.go",
        "listing_release.go",
        "marketplace.go",
    ],
    data = [
        "//:dotenv",
        "//data:ent_db",
        "//data:gorm_db",
    ],
    importpath = "github.com/m-nny/universe/lib/discsearch",
    visibility = ["//visibility:public"],
    deps = [
        "//lib/brain",
        "//lib/discogs",
        "//lib/spotify",
        "@com_github_joho_godotenv//:godotenv",
        "@com_github_mattn_go_sqlite3//:go-sqlite3",
    ],
)
