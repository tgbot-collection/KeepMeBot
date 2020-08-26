FROM docker

ENV GO111MODULE=on

WORKDIR /APP

RUN apk update && apk add --no-cache git make musl-dev go && \
git clone https://github.com/BennyThink/KeepMeBot /APP && go build -o main .

CMD /APP/main

# usage
# docker build -t keepmebot .
# docker run -d -e TOKEN="13faT8" keepmebot