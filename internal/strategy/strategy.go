package strategy

import (
	"sync"
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

type TradingStrategyDecisionEngine struct {
	Strategies []StrategyImplementation
}

func (engine *TradingStrategyDecisionEngine) GetAggregateDecisions(symbol string) AggregateResult {
	if engine.Strategies == nil || len(engine.Strategies) == 0 {
		panic("no strategies to apply in TradingStrategyEngine")
	}

	resultChannels := make([]<-chan StrategyResult, len(engine.Strategies))
	for i, strat := range engine.Strategies {
		resultChannels[i] = strat.ApplyStrategy(symbol)
	}

	result := AggregateResult{
		Decision:   Undecided.String(),
		Symbol:     symbol,
		Confidence: 0.0,
	}

	for resp := range processChannels(resultChannels...) {
		// todo: calculate the actual confidence and strategy across different strats
		result.Decision = resp.Decision.String()
		result.Date = resp.Date
		result.Confidence = 100
	}

	return result
}

func processChannels(channels ...<-chan StrategyResult) <-chan StrategyResult {
	var wg sync.WaitGroup

	wg.Add(len(channels))
	fanin := make(chan StrategyResult)
	multiplex := func(c <-chan StrategyResult) {
		defer wg.Done()
		for i := range c {
			fanin <- i
		}
	}
	for _, c := range channels {
		go multiplex(c)
	}
	go func() {
		wg.Wait()
		close(fanin)
	}()
	return fanin
}
