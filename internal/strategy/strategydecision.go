package strategy

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
