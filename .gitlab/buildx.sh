#!/bin/bash

args=(
    --file docker/Dockerfile
    --build-arg DOCKER_VERSION=${DOCKER_VERSION}
    --build-arg GOLANG_VERSION=${GOLANG_VERSION}
    --build-arg NODE_VERSION=${NODE_VERSION}
    --build-arg ALPINE_VERSION=${ALPINE_VERSION}
    --build-arg RCLONE_VERSION=${RCLONE_VERSION}
    --build-arg RESTIC_VERSION=${RESTIC_VERSION}
    --build-arg BUILDKIT_INLINE_CACHE=1
)

docker buildx create --use

docker pull --platform=linux/amd64 ${GO_CACHE} || true
docker pull --platform=linux/arm64 ${GO_CACHE} || true
docker buildx build . ${args[@]} \
    --target=goBuilder \
    --provenance=false \
    --platform=linux/amd64,linux/arm64 \
    --cache-from ${GO_CACHE} \
    --tag ${GO_CACHE} \
    --pull --push

docker pull --platform=linux/amd64 ${NODE_CACHE} || true
docker pull --platform=linux/arm64 ${NODE_CACHE} || true
docker buildx build . ${args[@]} \
    --target=nodeBuilder \
    --provenance=false \
    --platform=linux/amd64,linux/arm64 \
    --cache-from ${NODE_CACHE} \
    --tag ${NODE_CACHE} \
    --pull --push

docker pull --platform=linux/amd64 ${RESTIC_CACHE} || true
docker pull --platform=linux/arm64 ${RESTIC_CACHE} || true
docker buildx build . ${args[@]} \
    --target=resticBuilder \
    --provenance=false \
    --platform=linux/amd64,linux/arm64 \
    --cache-from ${RESTIC_CACHE} \
    --tag ${RESTIC_CACHE} \
    --pull --push

docker pull --platform=linux/amd64 ${LATEST_IMAGE} || true
docker pull --platform=linux/arm64 ${LATEST_IMAGE} || true
docker buildx build . ${args[@]} \
    --provenance=false \
    --platform=linux/amd64,linux/arm64 \
    --cache-from ${GO_CACHE} \
    --cache-from ${NODE_CACHE} \
    --cache-from ${RESTIC_CACHE} \
    --cache-from ${LATEST_IMAGE} \
    --build-arg APP_VERSION=${CI_COMMIT_TAG} \
    --build-arg BUILD_TIME=${CI_JOB_STARTED_AT} \
    --tag ${CURRENT_IMAGE} \
    --tag ${LATEST_IMAGE} \
    --pull --push
