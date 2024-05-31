package strategy

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoStrategies(t *testing.T) {
	engine := TradingStrategyDecisionEngine{
		Strategies: nil,
	}
	defer func() { recover() }()

	engine.GetAggregateDecisions("AAPL")

	// Never reaches here if `GetAggregateDecisions` panics.
	t.Errorf("did not panic")
}

type FakeStrategy struct {
	fakeResult StrategyResult
}

// returns whatever result I create the struct with
func (strat *FakeStrategy) ApplyStrategy(symbol string) <-chan StrategyResult {
	channel := make(chan StrategyResult)
	go func() {
		channel <- strat.fakeResult
		close(channel)
	}()
	return channel
}

func TestReturnsDecision(t *testing.T) {
	symbol := "AAPL"
	expected := StrategyResult{
		Symbol:   symbol,
		Decision: Buy,
		Success:  true,
	}
	fakeStrat := FakeStrategy{
		fakeResult: expected,
	}

	engine := TradingStrategyDecisionEngine{
		Strategies: []StrategyImplementation{
			&fakeStrat,
		},
	}

	expectedResp := AggregateResult{
		Decision:   Buy.String(),
		Symbol:     symbol,
		Confidence: 100,
	}
	resp := engine.GetAggregateDecisions(symbol)

	assert.Equal(t, expectedResp, resp)
}
