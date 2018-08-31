package main

import (
	"fmt"
	iex "github.com/jonwho/go-iex"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"strconv"
)

func formatQuote(quote *iex.Quote) string {
	stringOrder := []string{
		"Symbol",
		"Company Name",
		"Current",
		"High",
		"Low",
		"Open",
		"Close",
		"Change % (1 day)",
		"Volume",
	}

	outputMap := map[string]string{
		"Symbol":           quote.Symbol,
		"Company Name":     quote.CompanyName,
		"Current":          fmt.Sprintf("%#v", quote.LatestPrice),
		"High":             fmt.Sprintf("%#v", quote.High),
		"Low":              fmt.Sprintf("%#v", quote.Low),
		"Open":             fmt.Sprintf("%#v", quote.Open),
		"Close":            fmt.Sprintf("%#v", quote.Close),
		"Change % (1 day)": fmt.Sprintf("%#v", quote.ChangePercent) + " %",
		"Volume":           fmt.Sprintf("%#v", quote.LatestVolume),
	}

	printer := message.NewPrinter(language.English)
	fmtStr := "```\n"

	for _, k := range stringOrder {
		if k == "Volume" {
			i, _ := strconv.Atoi(outputMap[k])
			fmtStr += printer.Sprintf("%-17s %d\n", k, i)
		} else {
			fmtStr += printer.Sprintf("%-17s %-20s\n", k, outputMap[k])
		}
	}

	fmtStr += "```\n"

	return fmtStr
}
