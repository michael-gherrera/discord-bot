# discord-bot

## Get it running
* Get your bot token from discord from [here](https://discordapp.com/developers/applications/me).
* Create a `.env` file with `BOT_TOKEN=<your_token>`.
* Export the following variables for your system.
	- STOCK_API_URL
	- COIN_API_URL
	- ERR_INVALID_CMD
	- ERR_INVALID_SYM
* Install [golang/dep](https://github.com/golang/dep). On mac run `brew install dep`.
* Install go dependencies with `dep ensure`.
* Run `docker-compose up`
* Done!
