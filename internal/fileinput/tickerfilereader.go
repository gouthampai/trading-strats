package fileinput

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

type TickerProvider interface {
	GetTickers() []string
}
type TickerFileReader struct {
	FilePath string
	Logger   *log.Logger
}

func (tickerFileReader TickerFileReader) GetTickers() []string {
	jsonfile, err := os.Open(tickerFileReader.FilePath)
	if err != nil {
		tickerFileReader.Logger.Panic(err)
	}

	defer jsonfile.Close()

	byteValue, _ := io.ReadAll(jsonfile)

	var result []string
	err = json.Unmarshal([]byte(byteValue), &result)
	if err != nil {
		tickerFileReader.Logger.Panic(err)
	}

	return result
}
