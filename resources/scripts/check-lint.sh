#!/usr/bin/env bash

# detecting GOPATH and removing trailing "/" if any
GOPATH="$(go env GOPATH)"
GOPATH=${GOPATH%/}

# checking if golangci-lint is available
test -s "$GOPATH"/bin/golangci-lint || echo ">> installing golangci-lint" && GOBIN="$GOPATH"/bin go get -u github.com/golangci/golangci-lint/cmd/golangci-lint

this_path=$(dirname "$0")

echo ">> checking packages"
golangci-lint run -c "$this_path"/../lint/.golangci-pkg.yml ./cmd/... ./io/... ./version/... ./... || exit 1

if [[ -d "./internal" || -d "./cmd" ]]; then
    echo ">> checking unused exported symbols in internal packages"
    golangci-lint run -c "$this_path"/../lint/.golangci-internal.yml ./cmd/... || exit 1
fi
