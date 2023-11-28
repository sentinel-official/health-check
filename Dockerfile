FROM golang:1.21-alpine3.18 AS build

COPY . /root/

RUN --mount=target=/go/pkg/mod,type=cache \
    --mount=target=/root/.cache/go-build,type=cache \
    apk add make && \
    cd /root/ && make --jobs=$(nproc) install

FROM alpine:3.18

COPY --from=build /go/bin/02_client /usr/local/bin/main

RUN apk add --no-cache v2ray wireguard-tools && \
    rm -rf /etc/v2ray/ /usr/share/v2ray/

CMD ["main"]