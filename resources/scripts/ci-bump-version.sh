#!/usr/bin/env bash
set -e

# Configuration
git clone https://github.com/fsaintjacques/semver-tool /tmp/semver &> /dev/null
BUMPED_UP_VERSION=$(/tmp/semver/src/semver bump $SEMVER_RELEASE_LEVEL $VERSION)

echo "${BUMPED_UP_VERSION}"
