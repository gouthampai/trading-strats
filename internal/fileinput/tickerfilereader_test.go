package fileinput

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBadFilePath(t *testing.T) {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	reader := TickerFileReader{
		FilePath: "./aaa.json",
		Logger:   logger,
	}
	defer func() { recover() }()

	reader.GetTickers()

	// Never reaches here if `GetAggregateDecisions` panics.
	t.Errorf("did not panic")
}

func TestGoodJson(t *testing.T) {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	reader := TickerFileReader{
		FilePath: "./good.json",
		Logger:   logger,
	}

	result := reader.GetTickers()

	assert.NotNil(t, result)
	assert.Equal(t, 3, len(result))
}

func TestBadJson(t *testing.T) {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	reader := TickerFileReader{
		FilePath: "./bad.json",
		Logger:   logger,
	}

	defer func() { recover() }()

	reader.GetTickers()

	reader.GetTickers()

	// Never reaches here if `GetAggregateDecisions` panics.
	t.Errorf("did not panic")
}
