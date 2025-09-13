-- +goose Up
-- +goose StatementBegin
CREATE TABLE customers (
  id                     UUID PRIMARY KEY,
  user_id                UUID,
  email                  TEXT,
  default_payment_method JSONB,
  credit_balance_cents   BIGINT DEFAULT 0,
  created_at             TIMESTAMP WITH TIME ZONE DEFAULT now(),
  updated_at             TIMESTAMP WITH TIME ZONE DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE customers;
-- +goose StatementEnd
