package chart

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/BryanSLam/discord-bot/config"
	iex "github.com/jonwho/go-iex"
	chart "github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/seq"
	"github.com/wcharczuk/go-chart/util"
)

// Chart object that's used to create and manipulate any charts for a given ticker/symbol
type Chart struct {
	fileName string
}

// New returns a new chart object, takes in a filename to output charts to, otherwise uses default
func New(file string, ds iex.Chart) Chart {
	var name string
	if file != "" {
		if strings.Contains(file, ".") {
			if strings.Split(file, ".")[1] == "png" {
				name = file
			} else {
				panic(errors.New("Invalid file extension for chart"))
			}
		}
	} else {
		name = config.GetConfig().DefaultChartFileName
	}

	for _, data := range ds.Charts {
		parsed, err := time.Parse("15:04", data.Date)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}

		if data.Average != -1 {
			cd.XValues = append(cd.XValues, parsed)
			cd.YValues = append(cd.YValues, data.Average)
		}
	}

	return Chart{
		fileName: name,
	}
}

func (c *Chart) CreateChart(ticker string) {
	start := util.Date.Date(2018, 8, 17, util.Date.Eastern())
	end := util.Date.Date(2018, 8, 17, util.Date.Eastern())
	xv := seq.Time.MarketHours(start, end, util.NYSEOpen(), util.NYSEClose(), util.Date.IsNYSEHoliday)
	fmt.Println(d.YValues[len(d.YValues)-1])
	graph := chart.Chart{
		XAxis: chart.XAxis{
			Style:          chart.StyleShow(),
			TickPosition:   chart.TickPositionBetweenTicks,
			ValueFormatter: chart.TimeHourValueFormatter,
			Range: &chart.MarketHoursRange{
				MarketOpen:      util.NYSEOpen(),
				MarketClose:     util.NYSEClose(),
				HolidayProvider: util.Date.IsNYSEHoliday,
			},
		},
		YAxis: chart.YAxis{
			Style: chart.StyleShow(),
		},
		Series: []chart.Series{
			chart.TimeSeries{
				Name:    strings.ToUpper(ticker),
				XValues: xv,
				YValues: d.YValues,
			},
		},
	}

	fo, err := os.Create(c.fileName)
	if err != nil {
		panic(err)
	}
	err = graph.Render(chart.PNG, fo)

	if err != nil {
		panic(err)
	}
}
