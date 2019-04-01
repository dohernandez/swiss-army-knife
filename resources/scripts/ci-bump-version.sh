#!/usr/bin/env bash
set -e

# Find branch name of merge pull request
echo "Finding branch for commit ${TRAVIS_COMMIT}"
BRANCH_NAME=$(curl -H "Authorization: bearer ${GITHUB_TOKEN}" -X POST -d " \
 { \
   \"query\": \"query { \
  repository(owner:\\\"${GITHUB_OWNER}\\\", name:\\\"${GITHUB_REPO}\\\") { \
    pullRequests(states:MERGED, last: 10, orderBy: {field: UPDATED_AT, direction: ASC}){ \
      nodes{ \
        headRefName, \
        mergeCommit { \
          oid \
        } \
      } \
    } \
  } \
}\" \
 } \
 " https://api.github.com/graphql | jq -r ".data.repository.pullRequests.nodes[] | select(.mergeCommit.oid == \"${TRAVIS_COMMIT}\").headRefName")

echo "Found the following branch for commit ${TRAVIS_COMMIT}: ${BRANCH_NAME}"

# Find the largest version bump based on the merged PR's
BUMP=""

# Get the version bump based on the branch name
if echo "${BRANCH_NAME}" | grep -q -i -E '(^|[-/])(patch|issue|hotfix)[-/]?'; then
    BUMP='patch'
elif echo "${BRANCH_NAME}" | grep -q -i -E '(^|[-/])(minor|feature)[-/]?'; then
    BUMP='minor'
elif echo "${BRANCH_NAME}" | grep -q -i -E '(^|[-/])(major|release)[-/]?'; then
    BUMP='major'
else
    echo "Branch ${BRANCH_NAME}: Has a invalid branch name!"
    exit 1
fi

# Bump the version
echo "Bumping ${BUMP} version: ${VERSION}"

git clone https://github.com/fsaintjacques/semver-tool /tmp/semver &> /dev/null
BUMPED_UP_VERSION=$(/tmp/semver/src/semver bump $BUMP $VERSION)

echo "Bumped up ${BUMPED_UP_VERSION}"
echo "${BUMPED_UP_VERSION}" > version/version
