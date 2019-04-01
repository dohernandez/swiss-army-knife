#!/usr/bin/env bash
set -e

# Configuration
echo "Configuration"
VERSION=$(cat version/version)
echo "VERSION ${VERSION}"

ROOT_DIR="$(pwd)"
OUTPUT_DIR="${ROOT_DIR}/build"

# Installing vendor
make deps-vendor

# Build binaries
# shellcheck disable=SC2043
for OS in linux darwin windows; do
  for ARCH in amd64; do
    echo "Building binary for $OS/$ARCH..."
    BUILD_DIR="${OUTPUT_DIR}/${BINARY_NAME}_${VERSION}_${OS}_${ARCH}"

    # Build go binary
    GOARCH=${ARCH} GOOS=${OS} CGO_ENABLED=0 BUILD_DIR="${BUILD_DIR}" VERSION="${VERSION}" make build
    done
done

git tag "v${VERSION}"

# Archive binaries
cd "${OUTPUT_DIR}"

ARCHIVE_DIR="${ROOT_DIR}/archive"
mkdir "${ARCHIVE_DIR}"

# Package outputs
for i in ./*; do
    RELEASE=$(basename "${i}")

    echo "Packing binary for ${RELEASE}..."
    tar -czf "${ARCHIVE_DIR}/${RELEASE}.tar.gz" "${RELEASE}"
done
