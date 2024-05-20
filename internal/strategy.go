package internal

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
}

type AggregateResult struct {
	Decision   StrategyDecision
	Symbol     string
	Confidence float64
}

type StrategyImplementation interface {
	ApplyStrategy(symbol string) StrategyResult
}

type TradingStrategyDecisionEngine struct {
}

func (engine *TradingStrategyDecisionEngine) GetAggregateDecisions(symbol string) AggregateResult {
	//  todo: get all decisions from every strategy applied to a specific stock symbol
	// return an AggregateResult which specifies the majority decision returned by all strats and the confidence applied to the decision
}
