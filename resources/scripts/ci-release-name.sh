#!/usr/bin/env bash

# IMPORTANT!
# Do not print any log or any other information in this file, just the result.
# This script is used to set a variable value.

set -e

# Configuration
VERSION=$(cat version/version)

# Setting release name
RELEASE_NAME="v${VERSION}"

if [ "${TRAVIS_EVENT_TYPE}" == "pull_request" ]; then
    RELEASE_NAME="${RELEASE_NAME}-draft"
fi

echo "${RELEASE_NAME}"
