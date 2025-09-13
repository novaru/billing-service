# TODO

## Project Structure
```
billing-service/
├── cmd/
│   └── server/
│       └── main.go                 # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go              # Configuration management
│   ├── handler/
│   │   ├── checkout.go            # Checkout endpoints
│   │   ├── webhook.go             # Payment webhook handlers
│   │   ├── usage.go               # Usage reporting endpoints
│   │   ├── apikey.go              # API key management
│   │   └── invoice.go             # Invoice endpoints
│   ├── service/
│   │   ├── plan.go                # Plan service
│   │   ├── subscription.go        # Subscription service
│   │   ├── apikey.go              # API key service
│   │   ├── usage.go               # Usage service
│   │   ├── invoice.go             # Invoice service
│   │   └── payment.go             # Payment service
│   ├── repository/
│   │   ├── plan.go                # Plan repository wrapper
│   │   ├── customer.go            # Customer repository wrapper
│   │   ├── subscription.go        # Subscription repository wrapper
│   │   ├── apikey.go              # API key repository wrapper
│   │   ├── invoice.go             # Invoice repository wrapper
│   │   ├── transaction.go         # Transaction repository wrapper
│   │   └── usage.go               # Usage repository wrapper
│   ├── middleware/
│   │   ├── auth.go                # API key authentication
│   │   ├── logging.go             # Request logging
│   │   └── cors.go                # CORS handling
│   └── types/
│       ├── plan.go                # Plan types
│       ├── customer.go            # Customer types
│       ├── subscription.go        # Subscription types
│       ├── apikey.go              # API key types
│       ├── invoice.go             # Invoice types
│       ├── transaction.go         # Transaction types
│       └── usage.go               # Usage types
├── db/
│   ├── migrations/
│   │   ├── 001_create_plans.sql
│   │   ├── 002_create_customers.sql
│   │   ├── 003_create_subscriptions.sql
│   │   ├── 004_create_api_keys.sql
│   │   ├── 005_create_invoices.sql
│   │   ├── 006_create_transactions.sql
│   │   └── 007_create_usage_events.sql
│   ├── query/
│   │   ├── plans.sql              # Plan queries
│   │   ├── customers.sql          # Customer queries
│   │   ├── subscriptions.sql      # Subscription queries
│   │   ├── api_keys.sql           # API key queries
│   │   ├── invoices.sql           # Invoice queries
│   │   ├── transactions.sql       # Transaction queries
│   │   └── usage_events.sql       # Usage queries
│   └── generated/                 # Generated sqlc code
├── tests/
│   ├── integration/
│   │   └── api_test.go           # Integration tests
│   └── mocks/                    # Test mocks
├── pkg/
│   ├── logger/
│   │   └── logger.go             # Logging utilities
│   └── utils/
│       └── uuid.go               # UUID utilities
├── .env                          # Environment variables
├── .env.example                  # Example environment file
├── docker-compose.yml            # Local development setup
├── Dockerfile                    # Container definition
├── Makefile                      # Build and development commands
├── go.mod                        # Go modules
├── go.sum                        # Go modules checksum
├── sqlc.yaml                     # SQLC configuration
└── README.md                     # Project documentation
```


