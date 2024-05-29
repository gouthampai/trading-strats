package strategy

import (
	"errors"
	"log"
	"time"

	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	"github.com/shopspring/decimal"
)

type GetBarsClient interface {
	GetBars(symbol string, req marketdata.GetBarsRequest) ([]marketdata.Bar, error)
}

type smaResult struct {
	DayOfYear            time.Time       `json:"dayOfYear"`
	FiftyDayAverage      decimal.Decimal `json:"FiftyDayAverage"`
	TwoHundredDayAverage decimal.Decimal `json:"TwoHundredDayAverage`
}

type SmaCrossStrategy struct {
	Client GetBarsClient
	Logger *log.Logger
}

func (strat *SmaCrossStrategy) ApplyStrategy(symbol string) <-chan StrategyResult {
	response := make(chan StrategyResult)
	go func() {
		averages, error := strat.CalculateMovingAverages(symbol)
		if error != nil {
			strat.Logger.Println(error)
			response <- StrategyResult{
				Success:  false,
				Decision: Undecided,
				Symbol:   symbol,
			}
			close(response)
			return
		}

		// we need at least 2 records to check if there is a golden or death cross
		if len(averages) < 2 {
			strat.Logger.Println("fewer than 2 records")

			response <- StrategyResult{
				Success:  false,
				Decision: Undecided,
				Symbol:   symbol,
			}
			close(response)
			return
		}

		goldenCrossDetected := false
		deathCrossDetected := false

		for i := 1; i < len(averages); i++ {
			// todo: test golden cross and death cross detection logic
			prevRecord := averages[i-1]
			curRecord := averages[i]

			// golden cross detected
			if prevRecord.FiftyDayAverage.Compare(curRecord.TwoHundredDayAverage) == -1 && curRecord.FiftyDayAverage.Compare(curRecord.TwoHundredDayAverage) > -1 {
				goldenCrossDetected = true
				deathCrossDetected = false
			}

			if prevRecord.FiftyDayAverage.Compare(curRecord.TwoHundredDayAverage) == 1 && curRecord.FiftyDayAverage.Compare(curRecord.TwoHundredDayAverage) < 1 {
				goldenCrossDetected = false
				deathCrossDetected = true
			}
		}

		result := StrategyResult{
			Success:  true,
			Decision: Undecided,
			Symbol:   symbol,
		}

		if goldenCrossDetected {
			result.Decision = Buy
		} else if deathCrossDetected {
			result.Decision = Sell
		} else {
			result.Decision = Hold
		}

		response <- result
		close(response)
	}()

	return response
}

// assume this is calculating from the current day
// future state, pass in a date
// store historical data to reduce api calls?
func (strat *SmaCrossStrategy) CalculateMovingAverages(symbol string) ([]smaResult, error) {
	// get the last 214 days of data so that we can compute the moving average data for the last two weeks
	resp, err := strat.Client.GetBars(symbol, marketdata.GetBarsRequest{
		Start:     time.Now().Local().AddDate(0, 0, -365),
		TimeFrame: marketdata.NewTimeFrame(1, marketdata.Day),
	})

	if err != nil {
		return nil, err
	}

	if len(resp) < 201 {
		return nil, errors.New("fewer than 201 days of results returned by alpaca")
	}

	two_hundred_day_agg := decimal.Zero
	fifty_day_agg := decimal.Zero

	for i := 0; i < 200; i++ {
		curBar := resp[i]
		if i > 149 {
			fifty_day_agg = fifty_day_agg.Add(decimal.NewFromFloat(curBar.VWAP))
		}

		two_hundred_day_agg = two_hundred_day_agg.Add(decimal.NewFromFloat(curBar.VWAP))
	}

	decimal_two_hundred := decimal.NewFromInt(200)
	decimal_fifty := decimal.NewFromInt(50)

	result := []smaResult{
		{
			DayOfYear:            resp[199].Timestamp,
			FiftyDayAverage:      fifty_day_agg.Div(decimal_fifty),
			TwoHundredDayAverage: two_hundred_day_agg.Div(decimal_two_hundred),
		},
	}

	for i := 1; i < len(resp)-199; i++ {
		barToRemove := resp[i-1]
		barToAdd := resp[199+i]
		two_hundred_day_agg = two_hundred_day_agg.Sub(decimal.NewFromFloat(barToRemove.VWAP)).Add(decimal.NewFromFloat(barToAdd.VWAP))

		barToRemove = resp[149+i]
		fifty_day_agg = fifty_day_agg.Sub(decimal.NewFromFloat(barToRemove.VWAP)).Add(decimal.NewFromFloat(barToAdd.VWAP))

		temp := smaResult{
			DayOfYear:            resp[199+i].Timestamp,
			FiftyDayAverage:      fifty_day_agg.Div(decimal_fifty),
			TwoHundredDayAverage: two_hundred_day_agg.Div(decimal_two_hundred),
		}

		result = append(result, temp)
	}
	return result, nil
}
