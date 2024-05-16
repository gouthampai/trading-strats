package internal

import "github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"

type StrategyDecision int

const (
	Hold StrategyDecision = iota
	Buy
	Sell
)

func (s StrategyDecision) String() string {
	return [...]string{"Hold", "Buy", "Sell"}[s]
}

func (s StrategyDecision) EnumIndex() int {
	return int(s)
}

type StrategyResult struct {
	Decision StrategyDecision
	Symbol   string
}

type StrategyImplementation interface {
	ApplyStrategy(symbol string) StrategyResult
}

type TradingStrategyDecisionEngine struct {
	marketClient *marketdata.Client
}
