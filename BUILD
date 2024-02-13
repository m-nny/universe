load("@gazelle//:def.bzl", "gazelle")

# gazelle:build_file_name BUILD

gazelle(name = "gazelle")

filegroup(
    name = "dotenv",
    srcs = [".env"],
    visibility = ["//visibility:public"],
)
