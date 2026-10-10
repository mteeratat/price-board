package main

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/mteeratat/price-board/app/internal/store"
)

func TestPage(t *testing.T) {
	tests := []struct {
		name   string
		prices []store.Price
		want   string
	}{
		{"shows price", []store.Price{{Symbol: "BTCUSDT", Price: "82830.01", FetchedAt: time.Now()}}, "82830.01"},
		{"empty table", []store.Price{}, "No data yet"},
		{"escapes html", []store.Price{{Symbol: "<script>", Price: "1", FetchedAt: time.Now()}}, "&lt;script&gt;"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			if err := page.Execute(&buf, tt.prices); err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(buf.String(), tt.want) {
				t.Errorf("page does not contain %q", tt.want)
			}
		})
	}
}
