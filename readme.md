## To add new dependency
```bash
# Install deps
$ go get gorm.io/gorm
$ go get gorm.io/driver/sqlite

# Tidy go mod, so it marks deps as direct
$ go mod tidy

$ bazel mod tidy
```