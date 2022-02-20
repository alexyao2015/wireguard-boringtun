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
    && git clone https://git.zx2c4.com/wireguard-tools.git .

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
            openresolv

RUN set -x \
    && mkdir -p /config \
    && chmod -R 755 /config

COPY --from=s6downloader /s6downloader /
COPY --from=boringtun /boringtun /usr/bin
COPY --from=wgtools /wgtools/install /
COPY --from=rootfs /rootfs /

ENV \
    WG_QUICK_USERSPACE_IMPLEMENTATION=boringtun \
    WG_SUDO=1

ENTRYPOINT ["/init"]
