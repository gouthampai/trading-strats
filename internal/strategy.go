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

type StrategyImplementation interface {
	ApplyStrategy(symbol string) StrategyResult
}

type TradingStrategyDecisionEngine struct {
}
