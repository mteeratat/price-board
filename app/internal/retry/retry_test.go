package retry

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestDo(t *testing.T) {
	tests := []struct {
		name      string
		failTimes int // fn fails this many times, then succeeds
		wantCalls int
		wantErr   bool
	}{
		{"first try ok", 0, 1, false},
		{"ok after 2 fails", 2, 3, false},
		{"always fails", 99, 3, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calls := 0
			err := Do(context.Background(), 3, time.Millisecond, func() error {
				calls++
				if calls <= tt.failTimes {
					return errors.New("fail")
				}
				return nil
			})

			if (err != nil) != tt.wantErr {
				t.Errorf("err = %v, wantErr %v", err, tt.wantErr)
			}
			if calls != tt.wantCalls {
				t.Errorf("calls = %d, want %d", calls, tt.wantCalls)
			}
		})
	}
}
