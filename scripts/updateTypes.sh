#!/bin/sh

./scripts/swagger.sh init
cd web && yarn run types:openapi
