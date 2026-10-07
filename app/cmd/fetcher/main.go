package main

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"github.com/mteeratat/price-board/app/internal/binance"
	"github.com/mteeratat/price-board/app/internal/retry"
)

const (
	defaultPriceURL = "https://api.binance.com/api/v3/ticker/price"
	defaultSymbols  = "BTCUSDT,ETHUSDT"
)

func main() {
	url := getenv("BINANCE_PRICE_URL", defaultPriceURL)
	symbols := parseSymbols(getenv("SYMBOLS", defaultSymbols))
	if len(symbols) == 0 {
		log.Fatal("SYMBOLS is empty")
	}

	client := binance.NewClient(url)

	// Stop before the CronJob's activeDeadlineSeconds (50s) kills the pod.
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	failed := false
	for _, s := range symbols {
		var p *binance.Price
		err := retry.Do(ctx, 3, time.Second, func() error {
			var err error
			p, err = client.GetPrice(ctx, s)
			return err
		})
		if err != nil {
			log.Printf("fetch %s: %v", s, err)
			failed = true
			continue
		}
		log.Printf("%s = %s", p.Symbol, p.Price)
	}

	// Non-zero exit marks the Job as failed in Kubernetes.
	if failed {
		os.Exit(1)
	}
}

// getenv returns the env var, or fallback if it is unset or empty.
func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// parseSymbols turns "BTCUSDT, ETHUSDT," into ["BTCUSDT", "ETHUSDT"].
func parseSymbols(raw string) []string {
	var out []string
	for s := range strings.SplitSeq(raw, ",") {
		if s = strings.TrimSpace(s); s != "" {
			out = append(out, s)
		}
	}
	return out
}
