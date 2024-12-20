package main

import (
	"fmt"
	"time"

	"github.com/stevenwilkin/deribit-funding-transactions/deribit"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	d := deribit.NewDeribitFromEnv()

	for _, t := range d.GetTransactions() {
		if t.InstrumentName != "BTC-PERPETUAL" {
			continue
		}

		date := time.UnixMilli(t.Timestamp).Format(time.DateOnly)
		fmt.Printf("%s\t%9.6f\t%9.2f\t%7.2f\n", date, t.InterestPl, t.Price, t.InterestPl*t.Price)
	}
}
