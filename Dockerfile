FROM alpine

# TODO
RUN apk update && apk add wget && wget -O /APP https://github.com/BennyThink/xxx
WORKDIR /APP

CMD /APP/restblog

# usage
# docker build -t yyetsbot .
# docker run -d --restart=always -e TOKEN="TOKEN" bennythink/yyetsbot