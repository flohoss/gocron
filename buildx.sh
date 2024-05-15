#!/bin/bash

args=(
    --file docker/Dockerfile
    --build-arg V_DOCKER=${V_DOCKER}
    --build-arg V_GOLANG=${V_GOLANG}
    --build-arg V_NODE=${V_NODE}
    --build-arg V_ALPINE=${V_ALPINE}
    --build-arg V_RCLONE=${V_RCLONE}
    --build-arg V_RESTIC=${V_RESTIC}
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
    --provenance=false \
    --platform=linux/amd64,linux/arm64 \
    --cache-from=type=registry,ref=${GO_CACHE} \
    --cache-from=type=registry,ref=${NODE_CACHE} \
    --build-arg APP_VERSION=${CI_COMMIT_TAG} \
    --build-arg BUILD_TIME=${CI_JOB_STARTED_AT} \
    --tag ${CURRENT_IMAGE} \
    --tag ${LATEST_IMAGE} \
    --pull --push
