package strategy

import (
	"errors"
	"log"
	"os"
	"testing"

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
