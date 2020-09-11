FROM alpine

WORKDIR /APP

RUN apk update && apk add --no-cache wget && \
wget https://github.com/BennyThink/KeepMeBot/releases/latest/download/keepmebot-linux-amd64 -O /APP/main && \
chmod +x /APP/main

CMD /APP/main

# usage
# docker build -t keepmebot .
# docker run -d -e TOKEN="13faT8" keepmebot
