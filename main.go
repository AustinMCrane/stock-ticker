package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type CurrentQuote struct {
	Symbol      string  `json:"symbol"`
	Open        float64 `json:"open"`
	Close       float64 `json:"close"`
	LatestPrice float64 `json:"latestPrice"`
	CompanyName string  `json:"companyName"`
}

type StockSymbol struct {
	Symbol      string  `json:"symbol"`
	AverageCost float64 `json:"averageCost"`
}
type ConfigFile struct {
	StockSymbols []StockSymbol `json:"stocks"`
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

func readConfigFile(path string) (error, ConfigFile) {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return err, ConfigFile{}
	}

	var config ConfigFile

	json.Unmarshal(raw, &config)

	return nil, config
}

func printQuote(quote CurrentQuote, cost float64) {
	str := fmt.Sprintf("%s --- cost: %f --- price: %f --- gains: %f\n", quote.CompanyName, cost, quote.LatestPrice, quote.LatestPrice-cost)
	fmt.Printf(str)
}

func main() {
	_, file := readConfigFile("./config.json")
	var stockSymbols []StockSymbol
	for _, s := range file.StockSymbols {
		stockSymbols = append(stockSymbols, s)
	}
	for _, symbol := range stockSymbols {
		err, quote := getStockQuote(symbol.Symbol)
		if err != nil {
			fmt.Println("error gettings quote: ", err)
			os.Exit(1)
		}

		printQuote(quote, symbol.AverageCost)
	}
}
