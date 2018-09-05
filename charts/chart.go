package chart

import (
	"os"
	"strings"
	"time"

	"github.com/BryanSLam/discord-bot/config"
	iex "github.com/jonwho/go-iex"
	chart "github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/util"
)

func CreateChart(ticker string, ds iex.Chart) {
	xv, yv := processData(ds)

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
				YValues: c.YValues,
			},
		},
	}

	fo, err := os.Create(config.GetConfig().DefaultChartFileName)
	if err != nil {
		panic(err)
	}
	err = graph.Render(chart.PNG, fo)

	if err != nil {
		panic(err)
	}
}

func processData(ds iex.Chart) ([]time.Time, []float64) {
	var (
		xValues []time.Time
		yValues []float64
	)

	return xValues, yValues
}
