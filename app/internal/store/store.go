package store

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	pool *pgxpool.Pool
}

// New connects to Postgres and pings it, so a bad URL fails here, not on first use.
func New(ctx context.Context, databaseURL string) (*Store, error) {
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("store: connect: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("store: ping: %w", err)
	}
	return &Store{pool: pool}, nil
}

func (s *Store) Close() {
	s.pool.Close()
}

// InsertPrice saves one price. fetched_at is set by the DB default (now()).
func (s *Store) InsertPrice(ctx context.Context, symbol, price string) error {
	_, err := s.pool.Exec(ctx, "INSERT INTO prices (symbol, price) VALUES ($1, $2)", symbol, price)
	if err != nil {
		return fmt.Errorf("store: insert %s: %w", symbol, err)
	}
	return nil
}
