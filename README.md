# KeepMeBot
Keep me running bot

[![Build Status](https://travis-ci.org/BennyThink/KeepMeBot.svg?branch=master)](https://travis-ci.org/BennyThink/KeepMeBot)

## Support platform
* docker hub

## command
```
start - start using this bot
help - what can it do
list - list your service
add - add new service
history - is writter by whom?
ping - ping server
```

## run and deployment
It's strongly recommend to use docker to run this bot 
because we're about to running some system commands from untrusted sources.

```shell script
git clone https://github.com/BennyThink/KeepMeBot
cd KeepMeBot
# change your token here, you may also add other environment variables such as `http_proxy`
vim config.env
# create your db
touch keep.db
docker-compose up -d
```
Of course you could build your own docker image
`docker build -t keepmebot .

### How to update using docker-compose
1. Use docker pull to update docker-image, and run again

## License
MIT
