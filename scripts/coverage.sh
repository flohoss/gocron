#!/bin/sh
# Runs the Go test suite with coverage and strips generated files (sqlc, etc.)
# from the profile. Generated files are detected by the standard
# "// Code generated ... DO NOT EDIT." marker, so new generated code is
# excluded automatically without maintaining a path list.
set -eu

output="${1:-coverage.out}"
raw="${output}.raw"

go test ./... -covermode=atomic -coverprofile="$raw"

generated=$(mktemp)
trap 'rm -f "$raw" "$generated"' EXIT

grep -rl --include='*.go' -E '^// Code generated .+ DO NOT EDIT\.$' . |
	sed 's|^\./||; s|$|:|' >"$generated"

{
	head -1 "$raw"
	grep -vF -f "$generated" "$raw" | tail -n +2 || true
} >"$output"
