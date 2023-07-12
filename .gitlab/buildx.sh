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

docker buildx build . ${args[@]} \
    --target=goBuilder \
    --platform=linux/amd64,linux/arm64 \
    --cache-from=type=registry,ref=${GO_CACHE} \
    --cache-to=type=registry,ref=${GO_CACHE}

docker buildx build . ${args[@]} \
    --target=nodeBuilder \
    --platform=linux/amd64,linux/arm64 \
    --cache-from=type=registry,ref=${NODE_CACHE} \
    --cache-to=type=registry,ref=${NODE_CACHE}

docker buildx build . ${args[@]} \
    --target=resticBuilder \
    --provenance=false \
    --platform=linux/amd64,linux/arm64 \
    --cache-from=type=registry,ref=${RESTIC_CACHE} \
    --cache-to=type=registry,ref=${RESTIC_CACHE}

docker buildx build . ${args[@]} \
    --provenance=false \
    --platform=linux/amd64,linux/arm64 \
    --cache-from=type=registry,ref=${GO_CACHE} \
    --cache-from=type=registry,ref=${NODE_CACHE} \
    --cache-from=type=registry,ref=${RESTIC_CACHE} \
    --build-arg APP_VERSION=${CI_COMMIT_TAG} \
    --build-arg BUILD_TIME=${CI_JOB_STARTED_AT} \
    --tag ${CURRENT_IMAGE} \
    --tag ${LATEST_IMAGE} \
    --pull --push
