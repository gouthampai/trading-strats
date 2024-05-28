package strategy

import (
	"errors"
	"log"
	"time"

	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	"github.com/shopspring/decimal"
)

type smaResult struct {
	dayOfYear            time.Time
	fiftyDayAverage      decimal.Decimal
	twoHundredDayAverage decimal.Decimal
}

type SmaCrossStrategy struct {
	MarketClient *marketdata.Client
	Logger       *log.Logger
}

func (strat *SmaCrossStrategy) ApplyStrategy(symbol string) <-chan StrategyResult {
	response := make(chan StrategyResult)
	go func() {
		averages, error := strat.CalculateMovingAverages(symbol)
		if error != nil {
			strat.Logger.Fatal(error)
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
			// todo: implement golden cross and death cross detection logic
			prevRecord := averages[i-1]
			curRecord := averages[i]

			// golden cross detected
			if prevRecord.fiftyDayAverage.Compare(curRecord.twoHundredDayAverage) == -1 && curRecord.fiftyDayAverage.Compare(curRecord.twoHundredDayAverage) > -1 {
				goldenCrossDetected = true
				deathCrossDetected = false
			}

			if prevRecord.fiftyDayAverage.Compare(curRecord.twoHundredDayAverage) == 1 && curRecord.fiftyDayAverage.Compare(curRecord.twoHundredDayAverage) < 1 {
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
	resp, err := strat.MarketClient.GetBars(symbol, marketdata.GetBarsRequest{
		Start:     time.Now().Local().AddDate(0, 0, -365),
		TimeFrame: marketdata.NewTimeFrame(1, marketdata.Day),
	})

	if err != nil {
		return nil, err
	}

	if len(resp) < 200 {
		return nil, errors.New("Fewer than 200 days of results returned by alpaca")
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
			dayOfYear:            resp[200].Timestamp,
			fiftyDayAverage:      fifty_day_agg.Div(decimal_fifty),
			twoHundredDayAverage: two_hundred_day_agg.Div(decimal_two_hundred),
		},
	}

	for i := 1; i < len(resp)-200; i++ {
		barToRemove := resp[i-1]
		barToAdd := resp[199+i]
		two_hundred_day_agg = two_hundred_day_agg.Sub(decimal.NewFromFloat(barToRemove.VWAP)).Add(decimal.NewFromFloat(barToAdd.VWAP))

		barToRemove = resp[149+i]
		fifty_day_agg = fifty_day_agg.Sub(decimal.NewFromFloat(barToRemove.VWAP)).Add(decimal.NewFromFloat(barToAdd.VWAP))

		temp := smaResult{
			dayOfYear:            resp[200+i].Timestamp,
			fiftyDayAverage:      fifty_day_agg.Div(decimal_fifty),
			twoHundredDayAverage: two_hundred_day_agg.Div(decimal_two_hundred),
		}

		result = append(result, temp)
	}

	return result, nil
}
