CREATE TABLE prices (
  id         BIGSERIAL PRIMARY KEY,
  symbol     TEXT        NOT NULL,
  price      NUMERIC     NOT NULL,
  fetched_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_prices_symbol_time ON prices (symbol, fetched_at DESC);
