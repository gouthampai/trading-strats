package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
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
	flag.StringVar(&cfg.tickersJsonFilePath, "tickers-json-file-path", "/Users/gouthampai/Documents/code/trading-strats/cmd/bot/tickers.json", "Tickers json file path")
	flag.Parse()

	app := GenerateApplication(cfg)
	tickers := app.tickerProvider.GetTickers()
	processor := app.RegisterStrategyServices()
	processor.ProcessTickers(tickers)
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

func (app *application) RegisterStrategyServices() *strategy.TickerProcessor {
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

	processor := &strategy.TickerProcessor{
		Engine: *engine,
	}

	return processor
}
