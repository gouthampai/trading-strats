package strategy

import (
	"time"
)

type AggregateResult struct {
	Decision   string
	Symbol     string
	Confidence float64
	Date       time.Time
}
