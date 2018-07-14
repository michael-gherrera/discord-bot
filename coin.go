package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
)

type coin struct {
	Symbol        string
	Current       string
	Open          string
	High          string
	Low           string
	PercentChange string
	Response      string
}

func (c *coin) outputJSON() string {
	stringOrder := []string{
		"Symbol",
		"Current Price (USD)",
		"Open (24 Hours)",
		"High (24 Hours)",
		"Low (24 Hours)",
		"Change % (24 Hours)",
	}

	output := c.outputMap()

	fmtStr := "```\n"

	for _, k := range stringOrder {
		if k == "Change % (24 Hours)" {
			f, _ := strconv.ParseFloat(output[k], 64)
			fmtStr += fmt.Sprintf("%-20s %.2f %%\n", k, f)
		} else {
			fmtStr += fmt.Sprintf("%-20s %-20s\n", k, output[k])
		}
	}

	fmtStr += "```\n"

	return fmtStr
}

func (c *coin) UnmarshalJSON(data []byte) error {
	// auxiliary struct to help map json
	var aux struct {
		Display  map[string]interface{}
		Response string
	}

	dec := json.NewDecoder(bytes.NewReader(data))
	if err := dec.Decode(&aux); err != nil {
		return fmt.Errorf("decode coin: %v", err)
	}

	if aux.Response == "Error" {
		return fmt.Errorf("could not find coin: %v", c.Symbol)
	}

	c.Current = aux.Display[c.Symbol].(map[string]interface{})["USD"].(map[string]interface{})["PRICE"].(string)
	c.Open = aux.Display[c.Symbol].(map[string]interface{})["USD"].(map[string]interface{})["OPEN24HOUR"].(string)
	c.High = aux.Display[c.Symbol].(map[string]interface{})["USD"].(map[string]interface{})["HIGH24HOUR"].(string)
	c.Low = aux.Display[c.Symbol].(map[string]interface{})["USD"].(map[string]interface{})["LOW24HOUR"].(string)
	c.PercentChange = aux.Display[c.Symbol].(map[string]interface{})["USD"].(map[string]interface{})["CHANGEPCT24HOUR"].(string)

	return nil
}

func (c *coin) outputMap() map[string]string {
	return map[string]string{
		"Symbol":              c.Symbol,
		"Current Price (USD)": c.Current,
		"Open (24 Hours)":     c.Open,
		"High (24 Hours)":     c.High,
		"Low (24 Hours)":      c.Low,
		"Change % (24 Hours)": c.PercentChange,
	}
}
