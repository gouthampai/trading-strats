package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	"github.com/gouthampai/trading-strats/internal/strategy"
)

type config struct {
	port           int
	env            string
	alpacaEndpoint string
	alpacaKey      string
	alpacaSecret   string
}

type application struct {
	config        config
	logger        *log.Logger
	accountClient *alpaca.Client
	marketClient  *marketdata.Client
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "dev", "Environment (dev|stage|prod)")
	flag.StringVar(&cfg.alpacaEndpoint, "alpaca-endpoint", os.Getenv("ALPACA_ENDPOINT"), "Alpaca Endpoint")
	flag.StringVar(&cfg.alpacaKey, "alpaca-key", os.Getenv("ALPACA_KEY"), "Alpaca Key")
	flag.StringVar(&cfg.alpacaSecret, "alpaca-secret", os.Getenv("ALPACA_SECRET"), "Alpaca Secret")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	alpacaClient := alpaca.NewClient(alpaca.ClientOpts{
		APIKey:    cfg.alpacaKey,
		APISecret: cfg.alpacaSecret,
		BaseURL:   cfg.alpacaEndpoint,
	})

	marketClient := marketdata.NewClient(marketdata.ClientOpts{
		APIKey:    cfg.alpacaKey,
		APISecret: cfg.alpacaSecret,
	})

	app := &application{
		config:        cfg,
		logger:        logger,
		accountClient: alpacaClient,
		marketClient:  marketClient,
	}

	// todo: move this type creation into internal
	smaStrat := &strategy.SmaCrossStrategy{
		MarketClient: marketClient,
		Logger:       logger,
	}

	strats := []strategy.StrategyImplementation{
		smaStrat,
	}

	engine := &strategy.TradingStrategyDecisionEngine{
		Strategies: strats,
	}

	result := engine.GetAggregateDecisions("AMZN")

	app.prettyPrint(result)
}

func (app *application) prettyPrint(v any) {
	output, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		panic(err)
	}

	fmt.Printf(string(output))
}
