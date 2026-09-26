#!/bin/sh
# Generates a shields.io coverage badge SVG from a Go coverage profile.
# The color thresholds match common conventions so the badge visually
# reflects how close the project is to full coverage.
set -eu

profile="${1:-coverage.out}"
output="${2:-assets/coverage.svg}"

percent=$(go tool cover -func="$profile" | tail -1 | grep -Eo '[0-9]+\.[0-9]+')

color=$(awk -v p="$percent" 'BEGIN {
	if (p >= 90) print "brightgreen"
	else if (p >= 80) print "green"
	else if (p >= 70) print "yellowgreen"
	else if (p >= 60) print "yellow"
	else if (p >= 50) print "orange"
	else print "red"
}')

curl -sf "https://img.shields.io/badge/coverage-${percent}%25-${color}" -o "$output"

printf 'coverage: %s%% (%s)\n' "$percent" "$color"
