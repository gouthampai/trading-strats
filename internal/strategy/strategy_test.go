package strategy

import "testing"

func TestNoStrategies(t *testing.T) {
	engine := TradingStrategyDecisionEngine{
		Strategies: nil,
	}
	defer func() { recover() }()

	engine.GetAggregateDecisions("AAPL")

	// Never reaches here if `OtherFunctionThatPanics` panics.
	t.Errorf("did not panic")
}
