FROM golang:1.17-alpine as generator-builder
WORKDIR /generator

COPY generator/src .
RUN set -x \
    && go build -o generator

FROM golang:alpine as wg-go-builder
WORKDIR /wg-go

RUN set -x \
    && apk add \
        build-base \
        git \
    && git clone https://git.zx2c4.com/wireguard-go . \
    && git fetch --tags \
    && latestTag=$(git describe --tags `git rev-list --tags --max-count=1`) \
    && git checkout $latestTag

RUN set -x \
    && make PREFIX="/install" DESTDIR="/wg-go" install \
    && mv /wg-go/install/bin/wireguard-go .

##################
# Boringtun
##################
FROM rust:alpine as boringtun-builder
WORKDIR /buildtmp

RUN set -x \
    && apk add \
        build-base \
        git \
    && git clone https://github.com/cloudflare/boringtun.git .

RUN set -x \
    && cargo build --bin boringtun --target x86_64-unknown-linux-musl --release \
    && mv target/x86_64-unknown-linux-musl/release/boringtun .

WORKDIR /boringtun

RUN mv /buildtmp/boringtun .

###################
# S6 Overlay
###################
FROM alpine:latest as s6downloader
WORKDIR /s6downloader

RUN set -x \
    && S6_OVERLAY_VERSION=$(wget --no-check-certificate -qO - https://api.github.com/repos/just-containers/s6-overlay/releases/latest | awk '/tag_name/{print $4;exit}' FS='[""]') \
    && S6_OVERLAY_VERSION=${S6_OVERLAY_VERSION:1} \
    && wget -O /tmp/s6-overlay-arch.tar.xz "https://github.com/just-containers/s6-overlay/releases/download/v${S6_OVERLAY_VERSION}/s6-overlay-x86_64-${S6_OVERLAY_VERSION}.tar.xz" \
    && wget -O /tmp/s6-overlay-noarch.tar.xz "https://github.com/just-containers/s6-overlay/releases/download/v${S6_OVERLAY_VERSION}/s6-overlay-noarch-${S6_OVERLAY_VERSION}.tar.xz" \
    && mkdir -p /tmp/s6 \
    && tar -Jxvf /tmp/s6-overlay-noarch.tar.xz -C /tmp/s6 \
    && tar -Jxvf /tmp/s6-overlay-arch.tar.xz -C /tmp/s6 \
    && cp -r /tmp/s6/* .

###################
# Wireguard-tools
###################
FROM alpine:latest as wgtools
WORKDIR /wgtools

RUN set -x \
    && apk add \
        build-base \
        git \
        libmnl-dev \
    && git clone https://git.zx2c4.com/wireguard-tools.git . \
    && git fetch --tags \
    && latestTag=$(git describe --tags `git rev-list --tags --max-count=1`) \
    && git checkout $latestTag

WORKDIR /wgtools/src

RUN set -x \
    && make WITH_WGQUICK=yes PREFIX="/install" DESTDIR="/wgtools" install

###################
# Rootfs
###################
FROM alpine:latest as rootfs
WORKDIR /rootfs

COPY rootfs .

RUN chmod -R +x *

FROM alpine:latest

# wireguard-tools requires bash
# openresolv is required to set dns

RUN set -x \
    && apk add \
            bash \
            iptables \
            ip6tables \
            iproute2 \
            openresolv \
            libqrencode

RUN set -x \
    && mkdir -p \
          /config \
          /log \
          /run/wireguard \
    && chmod -R 755 \
          /config \
          /log \
    && chmod -R 600 \
          /run/wireguard \
    && chown -R nobody:nobody \
          /log

COPY --from=s6downloader /s6downloader /
COPY --from=boringtun-builder /boringtun /usr/bin
# COPY --from=wg-go-builder /wg-go/wireguard-go /usr/bin
COPY --from=wgtools /wgtools/install /
COPY --from=generator-builder /generator/generator /usr/bin
COPY --from=rootfs /rootfs /

ENV \
    S6_BEHAVIOUR_IF_STAGE2_FAILS=2 \
    S6_LOGGING_SCRIPT="n10 s1000000 T 1 T" \
    S6_CMD_WAIT_FOR_SERVICES_MAXTIME=0

ENV \
    WG_QUICK_USERSPACE_IMPLEMENTATION=boringtun \
    WG_SUDO=1

ENTRYPOINT ["/init"]
