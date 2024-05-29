package strategy

import (
	"errors"
	"log"
	"os"
	"testing"
	"time"

	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	"github.com/stretchr/testify/assert"
)

func TestSmaCrossStrategyReturnsErrorWhenInsufficientData(t *testing.T) {
	bars := make([]marketdata.Bar, 1)
	client := FakeClient{
		FakeBarsResult: bars,
		FakeError:      nil,
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	strat := SmaCrossStrategy{
		Client: &client,
		Logger: logger,
	}

	resp := strat.ApplyStrategy("AAPL")

	result := <-resp
	assert.Equal(t, false, result.Success)
	assert.Equal(t, "AAPL", result.Symbol)
	assert.Equal(t, Undecided, result.Decision)
}

func TestSmaCrossStrategyReturnsErrorWhenClientError(t *testing.T) {
	client := FakeClient{
		FakeBarsResult: nil,
		FakeError:      errors.New("some error from alpaca"),
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	strat := SmaCrossStrategy{
		Client: &client,
		Logger: logger,
	}

	resp := strat.ApplyStrategy("AAPL")

	result := <-resp
	assert.Equal(t, false, result.Success)
	assert.Equal(t, "AAPL", result.Symbol)
	assert.Equal(t, Undecided, result.Decision)
}

func TestCorrectDecisionIsReturnedByStrategy(t *testing.T) {
	tests := map[string]struct {
		decision StrategyDecision
	}{
		"hold": {
			decision: Hold,
		},
		"buy": {
			decision: Buy,
		},
		"sell": {
			decision: Sell,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			bars := GenerateBarsForTestData(test.decision)
			client := FakeClient{
				FakeBarsResult: bars,
				FakeError:      nil,
			}

			logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

			strat := SmaCrossStrategy{
				Client: &client,
				Logger: logger,
			}

			resp := strat.ApplyStrategy("AAPL")

			result := <-resp
			assert.Equal(t, true, result.Success)
			assert.Equal(t, "AAPL", result.Symbol)
			assert.Equal(t, test.decision, result.Decision)
		})
	}

}

type FakeClient struct {
	FakeBarsResult []marketdata.Bar
	FakeError      error
}

func (client *FakeClient) GetBars(symbol string, req marketdata.GetBarsRequest) ([]marketdata.Bar, error) {
	if client.FakeError != nil {
		return nil, client.FakeError
	} else {
		return client.FakeBarsResult, nil
	}
}

func GenerateBarsForTestData(desiredDecision StrategyDecision) []marketdata.Bar {
	barTimeStamp := time.Now().Local().AddDate(0, 0, -201)
	bars := make([]marketdata.Bar, 201)
	for i := 0; i < 201; i++ {
		if i < 200 {
			bars[i] = marketdata.Bar{
				Timestamp:  barTimeStamp,
				Open:       100,
				Close:      100,
				High:       100,
				Low:        100,
				Volume:     100_000,
				VWAP:       100,
				TradeCount: 1000,
			}
		} else {
			if desiredDecision == Hold {
				bars[i] = marketdata.Bar{
					Timestamp:  barTimeStamp,
					Open:       100,
					Close:      100,
					High:       100,
					Low:        100,
					Volume:     100_000,
					VWAP:       100,
					TradeCount: 1000,
				}
			} else if desiredDecision == Buy {
				// forms golden cross
				bars[i] = marketdata.Bar{
					Timestamp:  barTimeStamp,
					Open:       110,
					Close:      110,
					High:       110,
					Low:        110,
					Volume:     100_000,
					VWAP:       110,
					TradeCount: 1000,
				}
			} else if desiredDecision == Sell {
				// forms death cross
				bars[i] = marketdata.Bar{
					Timestamp:  barTimeStamp,
					Open:       90,
					Close:      90,
					High:       90,
					Low:        90,
					Volume:     100_000,
					VWAP:       90,
					TradeCount: 1000,
				}
			}
		}

		barTimeStamp = barTimeStamp.AddDate(0, 0, 1)
	}

	return bars
}
