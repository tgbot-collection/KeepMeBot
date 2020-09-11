FROM golang:alpine

ENV GO111MODULE=on

WORKDIR /APP

RUN apk update && apk add --no-cache alpine-sdk git make musl-dev sqlite && \
git clone https://github.com/BennyThink/KeepMeBot /APP && sh autogen.sh && go build -o main .

CMD /APP/main

# usage
# docker build -t keepmebot .
# docker run -d -e TOKEN="13faT8" keepmebot