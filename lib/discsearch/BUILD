load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "discsearch",
    srcs = [
        "discsearch.go",
        "marketplace.go",
        "master.go",
    ],
    data = [
        "//:dotenv",
        "//data:ent_db"
    ],
    importpath = "github.com/m-nny/universe/lib/discsearch",
    visibility = ["//visibility:public"],
    deps = [
        "//ent",
        "//lib/discogs",
        "//lib/spotify",
        "//lib/utils",
        "@com_github_joho_godotenv//:godotenv",
        "@com_github_mattn_go_sqlite3//:go-sqlite3",
    ],
)
