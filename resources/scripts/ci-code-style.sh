#!/usr/bin/env bash
set -e

# Configuration
echo "Setting configuration"
FILE_EXTENSIONS='\.go$'

# Detect the changed files
echo "Detecting the changed files"
git diff --name-only "${TRAVIS_COMMIT_RANGE}" | (grep -i -E "${FILE_EXTENSIONS}" || true) > changed_files.txt

printf "Change count: "
CHANGE_COUNT=$(wc -l < changed_files.txt)
if [ "${CHANGE_COUNT}" = "0" ]; then
    echo "No files affected. Skipping"
    exit 0
fi
echo "Affected files: ${CHANGE_COUNT}"

# Code style checker begin
printf "Checking golint: "
err_count=0
while IFS= read -r file; do
  if ! golint -set_exit_status "$file"; then
    err_count=$((err_count+1))
  fi
done < changed_files.txt

if [ $err_count -gt 0 ]; then
  exit 1
fi
echo "PASS"
echo

printf "Checking gofmt: "
# shellcheck disable=SC2002
ERRS=$(cat changed_files.txt | xargs gofmt -l 2>&1 || true)
if [ -n "${ERRS}" ]; then
    echo "FAIL - the following files need to be gofmt'ed:"
    for e in ${ERRS}; do
        echo "    $e"
    done
    echo
    exit 1
fi
echo "PASS"
echo

printf "Checking goimports: "
# shellcheck disable=SC2046
ERRS=$(goimports -l $(cat changed_files.txt) 2>&1 || true)
if [ -n "${ERRS}" ]; then
    echo "FAIL"
    echo "${ERRS}"
    echo
    exit 1
fi
echo "PASS"
echo
