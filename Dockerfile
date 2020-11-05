FROM golang:alpine as builder

ENV GO111MODULE=on

RUN apk update && apk add alpine-sdk tzdata git make musl-dev sqlite && \
git clone https://github.com/tgbot-collection/KeepMeBot /build && cd /build \
&& sh autogen.sh && go build -a -ldflags '-s -w' -o keep .


FROM alpine

RUN apk update && apk add git docker-cli
ENV TZ=Asia/Shanghai

COPY --from=builder /build/keep /keep
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

WORKDIR /

ENTRYPOINT ["/keep"]

# usage
# docker build -t keepmebot .
# docker run -d -e TOKEN="13faT8" keepmebot