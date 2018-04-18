#!/bin/sh
set -eux

mkdir -p target

VERSION=$(git describe --dirty --always)
LDFLAGS="-X main.Version=$VERSION"

GOOS=linux GOARCH=amd64 go build -ldflags "$LDFLAGS" -o target/collect-owners-linux-amd64
GOOS=linux GOARCH=arm go build -ldflags "$LDFLAGS" -o target/collect-owners-linux-arm
GOOS=linux GOARCH=arm64 go build -ldflags "$LDFLAGS" -o target/collect-owners-linux-arm64
GOOS=linux GOARCH=mips64 go build -ldflags "$LDFLAGS" -o target/collect-owners-linux-mips64
GOOS=darwin GOARCH=amd64 go build -ldflags "$LDFLAGS" -o target/collect-owners-darwin-amd64
GOOS=windows GOARCH=amd64 go build -ldflags "$LDFLAGS" -o target/collect-owners-windows-amd64
