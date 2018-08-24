# discord-bot

## Get it running
* Get your bot token from discord from [here](https://discordapp.com/developers/applications/me).
* Create a `.env` file with `BOT_TOKEN=<your_token>`.
* Install [golang/dep](https://github.com/golang/dep). On mac run `brew install dep`.
* Install go dependencies with `dep ensure`.
* Run `docker-compose up`
* Done!

## Development
* Do not modify the `Gopkg.lock` directly!
* Add new package with `dep ensure -add <path/to/pkg>` which will update `Gopkg.*` for you.
