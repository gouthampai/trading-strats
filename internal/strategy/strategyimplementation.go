package strategy

type StrategyImplementation interface {
	ApplyStrategy(symbol string) <-chan StrategyResult
}
