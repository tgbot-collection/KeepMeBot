# KeepMeBot
Keep me running bot

[![Build Status](https://travis-ci.org/BennyThink/KeepMeBot.svg?branch=master)](https://travis-ci.org/BennyThink/KeepMeBot)

## Support platform
* docker hub

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
### How to update using docker-compose
1. Use docker pull to update docker-image, and run again

## License
MIT