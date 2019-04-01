#!/usr/bin/env bash

set -e

# Installing goveralls
go get github.com/mattn/goveralls

# Installing vendor
make deps-vendor

# Running tests
make test

# Upload coverage results
# Example https://docs.coveralls.io/go
$GOPATH/bin/goveralls -coverprofile=overalls.coverprofile -service=travis-ci -repotoken $CODECOV_TOKEN
