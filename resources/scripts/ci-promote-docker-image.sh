#!/usr/bin/env bash

set -e

# Configuration
echo "Setting configuration"
echo "VERSION ${VERSION}"

if [ "${QUAY_USERNAME}" == "" ] || [ "${QUAY_PASSWORD}" == "" ]; then
    echo "No QUAY_USERNAME or QUAY_PASSWORD defined. Skipping"
    exit 0
fi

QUAY_REPO_SLUG="quay.io/${QUAY_USERNAME}/${BINARY_NAME}"

# Log in quay.io to pull the images
docker login -u="${QUAY_USERNAME}" -p="${QUAY_PASSWORD}" quay.io
docker build -t "${QUAY_REPO_SLUG}:${VERSION}" \
       --build-arg VERSION \
       --build-arg USER="${QUAY_USERNAME}" \
       . --cache-from "${QUAY_REPO_SLUG}:latest"
docker tag "${QUAY_REPO_SLUG}:${VERSION}" "${QUAY_REPO_SLUG}:latest"
docker push "${QUAY_REPO_SLUG}:${VERSION}"
docker push "${QUAY_REPO_SLUG}:latest"
