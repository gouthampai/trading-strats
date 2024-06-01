package strategy

import (
	"time"
)

type StrategyDecision int

const (
	Undecided StrategyDecision = iota
	Hold
	Buy
	Sell
)

func (s StrategyDecision) String() string {
	return [...]string{"Undecided", "Hold", "Buy", "Sell"}[s]
}

func (s StrategyDecision) EnumIndex() int {
	return int(s)
}

type StrategyResult struct {
	Success  bool
	Decision StrategyDecision
	Symbol   string
	Date     time.Time
}

type AggregateResult struct {
	Decision   string
	Symbol     string
	Confidence float64
	Date       time.Time
}

type StrategyImplementation interface {
	ApplyStrategy(symbol string) <-chan StrategyResult
}
