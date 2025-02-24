ARG V_DOCKER=27.3.1
ARG V_GOLANG=1.23
ARG V_NODE=lts
ARG V_ALPINE=3
ARG V_RCLONE=1
ARG V_RESTIC=0.17.3
FROM rclone/rclone:${V_RCLONE} AS rclone
FROM restic/restic:${V_RESTIC} AS restic
FROM node:${V_NODE}-alpine AS node-builder
WORKDIR /app

COPY ./web/package.json ./web/yarn.lock ./
RUN yarn install --frozen-lockfile

COPY ./web/ ./
RUN yarn build

FROM docker:${V_DOCKER}-cli AS final
RUN apk add --update zip tzdata borgbackup dumb-init && \
    rm -rf /tmp/* /var/tmp/* /usr/share/man /var/cache/apk/*

# rclone
COPY --from=rclone --chmod=0755 \
    /usr/local/bin/rclone /usr/bin/rclone

# restic
COPY --from=restic --chmod=0755 \
    /usr/bin/restic /usr/bin/restic

WORKDIR /app
COPY ./config/config.yml ./config/config.yml
COPY --from=node-builder /app/dist/ ./web/

# goreleaser
ARG TARGETARCH
ENV TARGETARCH=${TARGETARCH:-amd64}
COPY gocron ./gocron

ARG APP_VERSION
ENV APP_VERSION=$APP_VERSION
ARG BUILD_TIME
ENV BUILD_TIME=$BUILD_TIME

EXPOSE 8156

ENTRYPOINT ["dumb-init", "--"]
CMD ["/app/gocron"]
