#!/usr/bin/env bash
set -e

# Configuration
GIT_LAST_TAG=$(git describe --abbrev=0 --tags 2> /dev/null || echo "v0.1.0")
GIT_LAST_TAG=${GIT_LAST_TAG:1}

echo "${GIT_LAST_TAG}"
