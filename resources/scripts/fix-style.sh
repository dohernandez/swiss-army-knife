#!/usr/bin/env bash

echo ">> fixing imports and fmt"

SOURCES_TO_LINT=$(find . -name '*.go' -not -path "./vendor/*")

echo "$SOURCES_TO_LINT"| while read -r file; do
    line=$(awk '/^import \(/{printf "%s:%s", NR, $0}' < "$file" | head | cut -d':' -f1);
    if [[ $line -gt 0 ]]; then
        ((line++))
        import_out=$(awk "NR==$(( line )),/(^\\))/" < "$file" | sed -e '$ d')
        last_line=$((line + $(echo "$import_out" | wc -l)))
        ((last_line--))

        import=$(echo "$import_out" | awk 'NF')

        awk -v m=$line -v n=$last_line "m <= NR && NR <= n {next} {print}" "$file" > "${file}.new"
        mv "${file}.new" "$file"

        if [[ ! -z "${import// }" ]]; then
            perl -i -lpe "print '$import' if \$. == $line" "$file"
        fi

    fi
done

# We let goimports do the ordering
test -s "$(go env GOPATH)"/bin/goimports || go get -u golang.org/x/tools/cmd/goimports
# shellcheck disable=SC2086
"$(go env GOPATH)"/bin/goimports -l -w ${SOURCES_TO_LINT} &>/dev/null

# Simplify code
# shellcheck disable=SC2086
gofmt -s -w ${SOURCES_TO_LINT} &>/dev/null
