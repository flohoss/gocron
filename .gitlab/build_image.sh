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

docker pull ${GO_BUILDER_IMAGE} || true
docker build . ${args[@]} \
    --target goBuilder \
    --cache-from ${GO_BUILDER_IMAGE} \
    --tag ${GO_BUILDER_IMAGE}
docker push ${GO_BUILDER_IMAGE}

docker pull ${NODE_BUILDER_IMAGE} || true
docker build . ${args[@]} \
    --target nodeBuilder \
    --cache-from ${NODE_BUILDER_IMAGE} \
    --tag ${NODE_BUILDER_IMAGE}
docker push ${NODE_BUILDER_IMAGE}

docker pull ${LOGO_BUILDER_IMAGE} || true
docker build . ${args[@]} \
    --target logo \
    --cache-from ${LOGO_BUILDER_IMAGE} \
    --tag ${LOGO_BUILDER_IMAGE}
docker push ${LOGO_BUILDER_IMAGE}

docker pull ${LATEST_IMAGE} || true
docker build . ${args[@]} \
    --cache-from ${GO_BUILDER_IMAGE} \
    --cache-from ${NODE_BUILDER_IMAGE} \
    --cache-from ${LOGO_BUILDER_IMAGE} \
    --cache-from ${LATEST_IMAGE} \
    --build-arg VERSION=${CI_COMMIT_TAG} \
    --build-arg BUILDTIME=${CI_JOB_STARTED_AT} \
    --tag ${CURRENT_IMAGE} \
    --tag ${LATEST_IMAGE}
