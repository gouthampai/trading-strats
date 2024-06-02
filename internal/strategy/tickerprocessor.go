package strategy

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

type TickerProcessor struct {
	Engine TradingStrategyDecisionEngine
}

func (processor TickerProcessor) ProcessTickers(symbols []string) {
	numWorkers := 2 // the number of workers does not seem to matter really. wtf?
	startTime := time.Now()
	totalTickers := len(symbols)
	symbolsChan := make(chan string, totalTickers)
	resultsChan := make(chan AggregateResult, totalTickers)
	buys := make([]AggregateResult, 0)
	cutoffTime := time.Now().AddDate(0, 0, -30)
	var wg sync.WaitGroup
	wg.Add(numWorkers)
	// start workers
	for i := 0; i < numWorkers; i++ {
		go processor.ProcessTickerChannel(symbolsChan, resultsChan, &wg)
	}

	var resultsWg sync.WaitGroup
	resultsWg.Add(1)
	go func() {
		var counter = 0
		defer resultsWg.Done()
		for result := range resultsChan {
			counter++
			if counter%10 == 0 {
				fmt.Printf("Processed %v of %v\n", counter, totalTickers)
			}

			if counter == totalTickers {
				fmt.Println("Done processing!")
			}

			if result.Decision == "Buy" && result.Date.After(cutoffTime) {
				buys = append(buys, result)
			}
		}
	}()

	// Distribute jobs and wait for completion
	for _, s := range symbols {
		symbolsChan <- s
	}

	close(symbolsChan)
	wg.Wait()

	// Ensure all results are collected
	close(resultsChan)

	resultsWg.Wait()

	prettyPrint(buys)
	endTime := time.Now()
	fmt.Printf("%v symbols processed in %v", totalTickers, endTime.Sub(startTime))
}

func (processor TickerProcessor) ProcessTickerChannel(symbols <-chan string, results chan<- AggregateResult, wg *sync.WaitGroup) {
	defer wg.Done()

	for symbol := range symbols {
		results <- processor.Engine.GetAggregateDecisions(symbol)
	}
}

func prettyPrint(v any) {
	output, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		panic(err)
	}

	fmt.Printf(string(output))
}
