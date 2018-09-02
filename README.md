# discord-bot

## Requirements
* Go 1.11
    * Download from here: https://golang.org/dl/
    * Update Go using instructions here: https://gist.github.com/nikhita/432436d570b89cab172dcf2894465753
    * Verify with `go version`
* Discord Application with Bot support
* Docker (optional)
## Get it running
* Get your bot token from discord from [here](https://discordapp.com/developers/applications/me).
* Create a `.env` file with `BOT_TOKEN=<your_token>`.
* Install the dependencies
    * Run `go mod download` to install dependencies to your local cache.
    * Or run `go mod vendor` to install dependencies to a vendor folder in the project.
* Verify dependencies are installed with `go mod verify`
* Run `docker-compose up`
* Done!