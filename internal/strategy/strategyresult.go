package strategy

import "time"

type StrategyResult struct {
	Success  bool
	Decision StrategyDecision
	Symbol   string
	Date     time.Time
}
