FROM golang:1.17-alpine

RUN set -x \
    && apk add build-base

RUN set -x \
    && go get -v \
        golang.org/x/tools/gopls@latest \
        github.com/go-delve/delve/cmd/dlv@v1.8.0

RUN set -x \
    && apk add libqrencode
