#!/bin/sh

(cd tailwind && yarn install --frozen-lockfile)

air -c .air.toml
