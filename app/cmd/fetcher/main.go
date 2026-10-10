package main

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"github.com/mteeratat/price-board/app/internal/binance"
	"github.com/mteeratat/price-board/app/internal/env"
	"github.com/mteeratat/price-board/app/internal/retry"
	"github.com/mteeratat/price-board/app/internal/store"
)

const (
	defaultPriceURL = "https://api.binance.com/api/v3/ticker/price"
	defaultSymbols  = "BTCUSDT,ETHUSDT"
)

func main() {
	url := env.Get("BINANCE_PRICE_URL", defaultPriceURL)
	symbols := parseSymbols(env.Get("SYMBOLS", defaultSymbols))
	if len(symbols) == 0 {
		log.Fatal("SYMBOLS is empty")
	}
	// No default: a missing DB secret must fail loudly.
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	client := binance.NewClient(url)

	// Stop before the CronJob's activeDeadlineSeconds (50s) kills the pod.
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	st, err := store.New(ctx, dbURL)
	if err != nil {
		log.Fatal(err)
	}

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
		if err := st.InsertPrice(ctx, p.Symbol, p.Price); err != nil {
			log.Print(err)
			failed = true
			continue
		}
		log.Printf("saved %s = %s", p.Symbol, p.Price)
	}

	// Close explicitly: os.Exit below skips deferred calls.
	st.Close()

	// Non-zero exit marks the Job as failed in Kubernetes.
	if failed {
		os.Exit(1)
	}
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
