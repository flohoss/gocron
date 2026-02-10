#!/bin/sh
set -e

URL="${1}"

echo "Waiting for $URL to be ready..."

for i in $(seq 1 10); do
	if wget -qO- "$URL" >/dev/null 2>&1; then
		echo "✓ $URL is ready!"
		exit 0
	fi
	echo "Attempt $i/10 failed, retrying in 1s..."
	sleep 1
done

echo "✗ $URL did not become ready after 10 attempts"
exit 1
