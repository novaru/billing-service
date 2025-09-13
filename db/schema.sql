-- users
CREATE TABLE users (
    id            UUID PRIMARY KEY,
    email         TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    name          TEXT NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- plans
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

-- customers (users who pay, can be a team, company, or individuals)
CREATE TABLE customers (
  id                     UUID PRIMARY KEY,
  user_id                UUID, -- reference to your users table if needed
  email                  TEXT,
  default_payment_method JSONB,
  credit_balance_cents   BIGINT DEFAULT 0,
  created_at             TIMESTAMP WITH TIME ZONE DEFAULT now(),
  updated_at             TIMESTAMP WITH TIME ZONE DEFAULT now()
);

-- subscriptions
CREATE TABLE subscriptions (
  id                      UUID PRIMARY KEY,
  customer_id             UUID REFERENCES customers(id),
  plan_id                 UUID REFERENCES plans(id),
  status                  TEXT NOT NULL, -- "trialing","active","past_due","canceled"
  trial_ends_at           TIMESTAMP WITH TIME ZONE,
  current_period_start    TIMESTAMP WITH TIME ZONE,
  current_period_end      TIMESTAMP WITH TIME ZONE,
  cancel_at_period_end    BOOLEAN DEFAULT FALSE,
  gateway_subscription_id TEXT, -- e.g. stripe subscription id
  metadata                JSONB,
  created_at              TIMESTAMP WITH TIME ZONE DEFAULT now(),
  updated_at              TIMESTAMP WITH TIME ZONE DEFAULT now()
);

-- api keys
CREATE TABLE api_keys (
  id              UUID PRIMARY KEY,
  customer_id     UUID REFERENCES customers(id),
  api_key         TEXT UNIQUE NOT NULL, -- hashed identifier or token id (store hashed)
  hashed_key      TEXT NOT NULL,
  revoked         BOOLEAN DEFAULT FALSE,
  created_at      TIMESTAMP WITH TIME ZONE DEFAULT now(),
  last_used_at    TIMESTAMP WITH TIME ZONE,
  meta            JSONB
);

-- invoices
CREATE TABLE invoices (
  id                 UUID PRIMARY KEY,
  customer_id        UUID REFERENCES customers(id),
  subscription_id    UUID REFERENCES subscriptions(id),
  gateway_invoice_id TEXT,
  status             TEXT NOT NULL, -- draft, issued, paid, failed, void
  amount_cents       BIGINT NOT NULL,
  currency           TEXT NOT NULL DEFAULT 'USD',
  issued_at          TIMESTAMP WITH TIME ZONE,
  due_at             TIMESTAMP WITH TIME ZONE,
  paid_at            TIMESTAMP WITH TIME ZONE,
  pdf_url            TEXT,
  metadata           JSONB,
  created_at         TIMESTAMP WITH TIME ZONE DEFAULT now()
);

-- transactions
CREATE TABLE transactions (
  id                 UUID PRIMARY KEY,
  invoice_id         UUID REFERENCES invoices(id),
  customer_id        UUID REFERENCES customers(id),
  gateway            TEXT NOT NULL, -- stripe,xendit,midtrans
  gateway_payment_id TEXT,
  amount_cents       BIGINT NOT NULL,
  currency           TEXT NOT NULL,
  status             TEXT NOT NULL, -- succeeded, failed, refunded
  idempotency_key    TEXT,
  raw_response       JSONB,
  created_at         TIMESTAMP WITH TIME ZONE DEFAULT now()
);

-- usage events (append-only)
CREATE TABLE usage_events (
  id              UUID PRIMARY KEY,
  customer_id     UUID REFERENCES customers(id),
  api_key_id      UUID REFERENCES api_keys(id),
  metric          TEXT NOT NULL, -- "requests", "tokens", "bandwidth"
  quantity        BIGINT NOT NULL,
  cost_cents      BIGINT DEFAULT 0,
  reported_at     TIMESTAMP WITH TIME ZONE DEFAULT now(),
  processed       BOOLEAN DEFAULT FALSE,
  metadata        JSONB
);

-- aggregated monthly usage (rollups)
CREATE TABLE usage_aggregates (
  id                UUID PRIMARY KEY,
  customer_id       UUID REFERENCES customers(id),
  period_start      DATE NOT NULL,
  period_end        DATE NOT NULL,
  metric            TEXT NOT NULL,
  total_quantity    BIGINT DEFAULT 0,
  total_cost_cents  BIGINT DEFAULT 0,
  updated_at        TIMESTAMP WITH TIME ZONE DEFAULT now(),
  UNIQUE(customer_id, period_start, metric)
);

