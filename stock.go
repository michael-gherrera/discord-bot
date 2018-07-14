package main

import (
	"fmt"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"strconv"
)

type stock struct {
	Quote struct {
		Symbol        string  `json:"symbol"`
		CompanyName   string  `json:"companyName"`
		Current       float32 `json:"latestPrice"`
		High          float32 `json:"high"`
		Low           float32 `json:"low"`
		Open          float32 `json:"open"`
		Close         float32 `json:"close"`
		PercentChange float32 `json:"changePercent"`
		Volume        int32   `json:"latestVolume"`
	} `json:"quote"`
}

func (s *stock) outputJSON() string {
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

	output := s.outputMap()

	printer := message.NewPrinter(language.English)
	fmtStr := "```\n"

	for _, k := range stringOrder {
		if k == "Volume" {
			i, _ := strconv.Atoi(output[k])
			fmtStr += printer.Sprintf("%-17s %d\n", k, i)
		} else {
			fmtStr += printer.Sprintf("%-17s %-20s\n", k, output[k])
		}
	}

	fmtStr += "```\n"

	return fmtStr
}

func (s *stock) outputMap() map[string]string {
	return map[string]string{
		"Symbol":           s.Quote.Symbol,
		"Company Name":     s.Quote.CompanyName,
		"Current":          fmt.Sprintf("%#v", s.Quote.Current),
		"High":             fmt.Sprintf("%#v", s.Quote.High),
		"Low":              fmt.Sprintf("%#v", s.Quote.Low),
		"Open":             fmt.Sprintf("%#v", s.Quote.Open),
		"Close":            fmt.Sprintf("%#v", s.Quote.Close),
		"Change % (1 day)": fmt.Sprintf("%#v", s.Quote.PercentChange*100) + " %",
		"Volume":           fmt.Sprintf("%#v", s.Quote.Volume),
	}
}
