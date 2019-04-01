#!/usr/bin/env bash

# IMPORTANT!
# Do not print any log or any other information in this file, just the result.
# This script is used to set a variable value.

set -e

# Fetching the last tag
GIT_LAST_TAG=$(git describe --abbrev=0 --tags 2> /dev/null || echo "v0.1.0")
GIT_LAST_TAG=${GIT_LAST_TAG:1}

echo "${GIT_LAST_TAG}"
