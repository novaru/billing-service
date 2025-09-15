-- +goose Up
-- +goose StatementBegin
CREATE TABLE plans (
  id            UUID PRIMARY KEY,
  slug          TEXT UNIQUE NOT NULL, -- "free", "starter", "pro"
  name          TEXT NOT NULL,
  description   TEXT,
  price_cents   BIGINT NOT NULL, -- price per period
  currency      TEXT NOT NULL DEFAULT 'USD',
  interval      TEXT NOT NULL, -- "month", "year"
  quota_limits  JSONB, -- {"requests_per_month": 100000, "tokens_per_month": 1000000}
  meta          JSONB,
  created_at    TIMESTAMP WITH TIME ZONE DEFAULT now(),
  updated_at    TIMESTAMP WITH TIME ZONE DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE plans;
-- +goose StatementEnd
