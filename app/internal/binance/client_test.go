package binance

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetPrice(t *testing.T) {
	tests := []struct {
		name    string
		symbol  string
		status  int
		body    string
		wantErr bool
	}{
		{"ok", "BTCUSDT", 200, `{"symbol":"BTCUSDT","price":"1"}`, false},
		{"server error", "BTCUSDT", 500, `{}`, true},
		{"bad json", "BTCUSDT", 200, `not json`, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if got := r.URL.Query().Get("symbol"); got != tt.symbol {
					t.Errorf("server got symbol %q, want %q", got, tt.symbol)
				}
				w.WriteHeader(tt.status)
				w.Write([]byte(tt.body))
			}))
			defer srv.Close()

			c := &Client{URL: srv.URL, HTTP: srv.Client()}
			_, err := c.GetPrice(context.Background(), tt.symbol)

			if (err != nil) != tt.wantErr {
				t.Errorf("err = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
