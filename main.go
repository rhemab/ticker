package main

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly"
)

type Stock struct {
	ticker string
	price  string
}

const baseUrl = "https://finance.yahoo.com/quote/"

var elapsedTime time.Duration

func main() {
	args := os.Args
	if len(args) == 1 {
		fmt.Println("Add tickers as arguments")
		return
	}
	ch := make(chan Stock, len(args))
	var wg sync.WaitGroup
	startTime := time.Now()

	for _, ticker := range args[1:] {
		wg.Add(1)
		go getTicker(&wg, strings.ToUpper(ticker), ch)
	}
	go func() {
		wg.Wait()
		close(ch)
		elapsedTime = time.Since(startTime)
	}()

	for stock := range ch {
		fmt.Println(stock.ticker, stock.price)
	}
	fmt.Println("Time Taken: ", elapsedTime)
}

func getTicker(wg *sync.WaitGroup, ticker string, ch chan<- Stock) {
	defer wg.Done()

	var stock Stock
	url := baseUrl + ticker

	collector := colly.NewCollector()

	collector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})
	collector.OnError(func(r *colly.Response, e error) {
		fmt.Println("Error:", e, ticker)
	})
	collector.OnHTML(".livePrice", func(e *colly.HTMLElement) {
		stock.ticker = ticker
		stock.price = e.ChildText("span")
	})
	collector.Visit(url)
	ch <- stock
}
