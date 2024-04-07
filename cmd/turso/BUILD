load("@rules_go//go:def.bzl", "go_binary", "go_library")

go_library(
    name = "turso_lib",
    srcs = ["turso.go"],
    importpath = "github.com/m-nny/universe/cmd/turso",
    visibility = ["//visibility:private"],
    deps = [
        "//lib/brain",
        "@com_github_jmoiron_sqlx//:sqlx",
        "@com_github_joho_godotenv//:godotenv",
        "@com_github_tursodatabase_libsql_client_go//libsql",
    ],
)

go_binary(
    name = "turso",
    embed = [":turso_lib"],
    visibility = ["//visibility:public"],
)
