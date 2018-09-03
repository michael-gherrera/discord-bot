package config

import (
	"fmt"

	"github.com/tkanos/gonfig"
)

type BotConfig struct {
	StockAPIURL           string
	CoinAPIURL            string
	InvalidCommandMessage string
	InvalidSymbolMessage  string
	DefaultChartFileName  string
}

var (
	config BotConfig
)

func GetConfig() *BotConfig {
	// Initalize config if it doesn't already exist
	if config == (BotConfig{}) {
		// Use gonfig to fetch the config variables from config.json
		err := gonfig.GetConf("config.json", &config)
		if err != nil {
			fmt.Println("error fetching config values", err)
			return nil
		}
	}

	return &config
}
