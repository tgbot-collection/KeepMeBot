FROM golang:alpine as builder

ENV GO111MODULE=on


RUN apk update && apk add alpine-sdk git make musl-dev sqlite && \
git clone https://github.com/tgbot-collection/KeepMeBot /build && cd /build \
&& sh autogen.sh && CGO_ENABLED=0 go build -a -ldflags '-s -w -extldflags "-static"' -o keep .


FROM scratch

COPY --from=builder /build/keep /keep
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

WORKDIR /

ENTRYPOINT ["/keep"]

# usage
# docker build -t keepmebot .
# docker run -d -e TOKEN="13faT8" keepmebot