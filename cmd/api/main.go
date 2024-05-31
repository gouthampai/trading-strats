package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	"github.com/gouthampai/trading-strats/internal/fileinput"
	"github.com/gouthampai/trading-strats/internal/strategy"
)

type config struct {
	port                int
	env                 string
	alpacaEndpoint      string
	alpacaKey           string
	alpacaSecret        string
	tickersJsonFilePath string
}

type application struct {
	config         config
	logger         *log.Logger
	accountClient  *alpaca.Client
	marketClient   *marketdata.Client
	tickerProvider fileinput.TickerProvider
}

func main() {
	fmt.Println("program starting at: " + time.Now().String())
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "dev", "Environment (dev|stage|prod)")
	flag.StringVar(&cfg.alpacaEndpoint, "alpaca-endpoint", os.Getenv("ALPACA_ENDPOINT"), "Alpaca Endpoint")
	flag.StringVar(&cfg.alpacaKey, "alpaca-key", os.Getenv("ALPACA_KEY"), "Alpaca Key")
	flag.StringVar(&cfg.alpacaSecret, "alpaca-secret", os.Getenv("ALPACA_SECRET"), "Alpaca Secret")
	flag.StringVar(&cfg.tickersJsonFilePath, "tickers-json-file-path", "/Users/gouthampai/Documents/code/trading-strats/cmd/api/tickers.json", "Tickers json file path")
	flag.Parse()

	app := GenerateApplication(cfg)
	tickers := app.tickerProvider.GetTickers()
	engine := app.RegisterStrategyServices()

	buys := make([]strategy.AggregateResult, 0)
	cutoffTime := time.Now().AddDate(0, 0, -30)
	// todo: refactor with a worker queue system to prevent 429 errors from alpaca.
	var wg sync.WaitGroup
	for _, symbol := range tickers {
		wg.Add(1)
		go func() {

			randomSeconds := rand.Int63n(100)
			time.Sleep(time.Second * time.Duration(randomSeconds))
			result := engine.GetAggregateDecisions(symbol)
			fmt.Println(time.Now())
			app.prettyPrint(result)
			if result.Decision == "Buy" && result.Date.After(cutoffTime) {
				buys = append(buys, result)
			}
			defer wg.Done()
		}()

	}

	wg.Wait()
	//app.prettyPrint(buys)
}

func GenerateApplication(config config) *application {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	alpacaClient := alpaca.NewClient(alpaca.ClientOpts{
		APIKey:    config.alpacaKey,
		APISecret: config.alpacaSecret,
		BaseURL:   config.alpacaEndpoint,
	})

	marketClient := marketdata.NewClient(marketdata.ClientOpts{
		APIKey:    config.alpacaKey,
		APISecret: config.alpacaSecret,
	})

	tickerProvider := fileinput.TickerFileReader{
		FilePath: config.tickersJsonFilePath,
		Logger:   logger,
	}

	app := &application{
		config:         config,
		logger:         logger,
		accountClient:  alpacaClient,
		marketClient:   marketClient,
		tickerProvider: tickerProvider,
	}

	return app
}

func (app *application) prettyPrint(v any) {
	output, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		panic(err)
	}

	fmt.Printf(string(output))
}

func (app *application) RegisterStrategyServices() *strategy.TradingStrategyDecisionEngine {
	smaStrat := &strategy.SmaCrossStrategy{
		Client: app.marketClient,
		Logger: app.logger,
	}

	strats := []strategy.StrategyImplementation{
		smaStrat,
	}

	engine := &strategy.TradingStrategyDecisionEngine{
		Strategies: strats,
	}

	return engine
}
