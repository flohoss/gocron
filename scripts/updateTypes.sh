#!/bin/sh

docker compose stop backend
./scripts/swagger.sh format
./scripts/swagger.sh init
docker compose start backend
cd web && yarn run types:openapi
