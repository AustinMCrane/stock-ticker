package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type CurrentQuote struct {
	Symbol      string  `json:"symbol"`
	Open        float64 `json:"open"`
	Close       float64 `json:"close"`
	LatestPrice float64 `json:"latestPrice"`
}

func getStockQuote(symbol string) (error, CurrentQuote) {
	url := fmt.Sprintf("https://api.iextrading.com/api/1.0/stock/%s/quote?displayPercent=true", symbol)
	resp, err := http.Get(url)
	if err != nil {
		return err, CurrentQuote{}
	}
	decoder := json.NewDecoder(resp.Body)
	var quote CurrentQuote
	err = decoder.Decode(&quote)
	if err != nil {
		return err, CurrentQuote{}
	}

	return nil, quote
}

func printQuote(quote CurrentQuote) {
	str := fmt.Sprintf("Stock: %s --- price: %f\n", quote.Symbol, quote.LatestPrice)
	fmt.Printf(str)
}

func main() {
	stockSymbols := os.Args[1:]
	for _, symbol := range stockSymbols {
		err, quote := getStockQuote(symbol)
		if err != nil {
			fmt.Println("error gettings quote: ", err)
			os.Exit(1)
		}

		printQuote(quote)
	}
}
