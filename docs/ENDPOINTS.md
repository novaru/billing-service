# API Endpoints Documentation

## Public Endpoints (No Authentication Required)

#### Authentication & User Management
`POST /api/v1/auth/register`<br>
Creates a new customer account<br>
Body: 
```js
{ name: "John Doe", email: "john@example.com", password: "password123" }
```
Response: 
```js 
{ 
  success: true, 
  data: { 
    id: "uuid", 
    email: "john@example.com",
    created_at: "..." 
  } 
}
```

`POST /api/v1/auth/login`<br>
Authenticates user and returns JWT token<br>
Body: 
```js 
{ email: "john@example.com", password: "password123" }
```
Response: 
```js 
{ 
  success: true, 
  data: { 
    token: "jwt_token", 
    expires_at: "..." 
  } 
}
```

`POST /api/v1/auth/forgot-password`<br>
Sends password reset email to user<br>
Body: 
```js 
{ email: "john@example.com" }
```
Response: 
```js
{ success: true, data: "Password reset email sent" }
```

`POST /api/v1/auth/reset-password`<br>
Resets user password using token from email<br>
Body: 
```js
{ token: "reset_token", new_password: "newpassword123" }
```
Response: 
```js
{ success: true, data: "Password reset successfully" }
```

#### Plans & Pricing
`GET /api/v1/plans`<br>
Lists all available subscription plans<br>
Query params: `?active=true` (optional, filter active plans only)<br>
Response: 
```js 
{ success: true, 
  data: [
    { 
      id: "uuid",
      name: "Pro", 
      price_cents: 2999,
      features: { /* ... */ } 
    }
  ]
}
```

`GET /api/v1/plans/{id}`<br>
Gets details of a specific plan<br>
Response: 
```js
{ 
  success: true, 
  data: { 
    id: "uuid",
    name: "Pro", 
    description: "...", 
    features: { /* ... */ } 
  } 
}
```

#### Checkout & Subscription
`POST /api/v1/checkout/create`<br>
Creates a checkout session for subscription (Stripe/payment gateway)<br>
Body: 
```js
{ plan_id: "uuid", customer_email: "john@example.com", success_url: "...", "cancel_url": "..." }
```
Response:
```js
{ 
  success: true, 
  data: { 
    checkout_url: "https://checkout.stripe.com/...", 
    session_id: "cs_..." 
    } 
}
```

`GET /api/v1/checkout/success`<br>
Handles successful checkout redirect (optional, could be handled on frontend)<br>
Query params: `?session_id=cs_...`<br>
Response: Redirect to success page or JSON confirmation


## Protected Endpoints (JWT Authentication Required)

#### User Profile Management
`GET /api/v1/profile`<br>
Gets current user's profile information<br>
Headers: `Authorization: Bearer <jwt_token>`<br>
Response:
```js
{
  success: true,
  data: {
    id: "uuid",
    name: "John",
    email: "john@example.com",
    created_at: "..."
  }
}
```

`PUT /api/v1/profile`<br>
Updates current user's profile information<br>
Headers: `Authorization: Bearer <jwt_token>`<br>
Body: 
```js
{ name: "John Updated", company: "Acme Inc" }
```
Response:
```js 
{ 
  success: true,
  data: {
    id: "uuid",
    name: "John Updated",
    /* ... */
  } 
}
```

`DELETE /api/v1/profile`<br>
Deletes current user's account (soft delete, cancel subscriptions)<br>
Headers: `Authorization: Bearer <jwt_token>`<br>
Response: 
```js
{ success: true, data: "Account deleted successfully" }
```


#### Subscription Management
`GET /api/v1/subscriptions`<br>
Gets current user's active subscription<br>
Headers: `Authorization: Bearer <jwt_token>`<br>
Response: 
```js
{ 
  success: true, 
  data: {
    id: "uuid",
    plan: { /* ... */ },
    status: "active",
    current_period_end: "..." 
  }
}
```

`POST /api/v1/subscriptions/cancel`<br>
Cancels current subscription (cancel at period end)<br>
Headers: `Authorization: Bearer <jwt_token>`<br>
Body: 
```js
{ reason: "Too expensive", cancel_immediately: false }
```
Response: 
```js
{ success: true, data: { canceled_at: "...", cancel_at_period_end: true } }
```

`POST /api/v1/subscriptions/reactivate`<br>
Reactivates a canceled subscription (before period ends)<br>
Headers: `Authorization: Bearer <jwt_token>`<br>
Response: 
```js 
{ 
  success: true,
  data: { 
    status: "active",
    cancel_at_period_end: false 
  }
}
```

`PUT /api/v1/subscriptions/upgrade`<br>
Upgrades/downgrades subscription plan<br>
Headers: `Authorization: Bearer <jwt_token>`<br>
Body: 
```js 
{ new_plan_id: "uuid", prorate: true }
```
Response:
```js
{ success: true, data: { subscription: { /* ... */ }, invoice: { /* ... */ } } }
```

`PUT /api/v1/subscriptions/payment-method`<br>
Updates payment method for subscription<br>
Headers: `Authorization: Bearer <jwt_token>`<br>
Body: 
```js 
{ payment_method_id: "pm_stripe_id" }
```
Response: 
```js
{ success: true, data: "Payment method updated successfully" }
```


#### API Key Management
`GET /api/v1/api-keys`<br>
Lists user's API keys (shows prefix only, not full key)<br>
Headers: `Authorization: Bearer <jwt_token>`<br>
Response: 
```js
{ 
  success: true,
  data: [{ 
    id: "uuid", 
    name: "Production Key",
    prefix: "sk_live_abc",
    created_at: "..." 
  }] 
}
```

`POST /api/v1/api-keys`<br>
Creates a new API key for the user<br>
Headers: `Authorization: Bearer <jwt_token>`<br>
Body: 
```js
{ name: "Production Key", scopes: ["usage:write", "usage:read"] }
```
Response: 
```js
{ success: true, data: { id: "uuid", key: "sk_live_abc123...", name: "Production Key" } }
```
NOTE: Full API key is only returned once during creation

`PUT /api/v1/api-keys/{id}`<br>
Updates API key name or scopes<br>
Headers: `Authorization: Bearer <jwt_token>`<br>
Body: 
```js
{ name: "Updated Name", scopes: ["usage:write"] }
```
Response: 
```js
{
  success: true,
  data: { 
    id: "uuid",
    name: "Updated Name", 
    /* ... */ 
  } 
}
```

`DELETE /api/v1/api-keys/{id}`<br>
Revokes/deletes an API key<br>
Headers: `Authorization: Bearer <jwt_token>`<br>
Response: 
```js
{ success: true, data: "API key revoked successfully" }
```

`POST /api/v1/api-keys/{id}/rotate`<br>
Rotates an API key (generates new key, invalidates old one)<br>
Headers: `Authorization: Bearer <jwt_token>`<br>
Response: 
```js 
{ 
  success: true,
  data: {
    id: "uuid",
    key: "sk_live_new123...",
    rotated_at: "..."
  } 
}
```


#### Usage & Analytics
`GET /api/v1/usage`<br>
Gets current usage statistics for the user<br>
Headers: `Authorization: Bearer <jwt_token>`<br>
Query params: `?period=current_month&group_by=day (optional)`<br>
Response: 
```js
{ 
  success: true,
  data: {
    tokens_used: 15000,
    requests_made: 500,
    quota: {
      tokens: 50000, 
      requests: 1000 
    } 
  }
}
```

`GET /api/v1/usage/history`<br>
Gets historical usage data<br>
Headers: `Authorization: Bearer <jwt_token>`<br>
Query params: `?start_date=2024-01-01&end_date=2024-01-31&granularity=daily`
Response:
```js
{ 
  success: true,
  data: [{
    date: "2024-01-01",
    tokens_used: 1000,
    requests_made: 50 
  }]
}
```

`GET /api/v1/usage/limits`<br>
Gets current usage limits and quotas<br>
Headers: `Authorization: Bearer <jwt_token>`<br>
Response: 
```js
{ 
  success: true,
  data: {
    plan: "Pro",
    limits: {
      tokens_monthly: 50000,
      requests_monthly: 1000 
    } 
  } 
}
```

#### Invoicing & Billing
`GET /api/v1/invoices`<br>
Lists user's invoices<br>
Headers: `Authorization: Bearer <jwt_token>`<br>
Query params: `?status=paid&limit=10&offset=0`
Response: 
```js 
{ 
  success: true,
  data: [{ 
    id: "uuid", 
    invoice_number: "INV-2024-001", 
    total_cents: 2999, 
    status: "paid" 
  }]
}
```

`GET /api/v1/invoices/{id}`<br>
Gets detailed invoice information<br>
Headers: `Authorization: Bearer <jwt_token>`<br>
Response: 
```js
{ 
  success: true, 
  data: {
    id: "uuid",
    line_items: [ /* ... */ ],
    total_cents: 2999,
    pdf_url: "..." 
  } 
}
```

`GET /api/v1/invoices/{id}/pdf`<br>
Downloads invoice as PDF<br>
Headers: `Authorization: Bearer <jwt_token>`<br>
Response: PDF file download

`POST /api/v1/invoices/{id}/pay`<br>
Manually pay an unpaid invoice (retry payment)<br>
Headers: `Authorization: Bearer <jwt_token>`<br>
Response: 
```js 
{ success: true, data: { "payment_intent_url": "https://..." } }
```


## API Key Protected Endpoints (For Service-to-Service)

#### Usage Reporting (Primary API)
`POST /api/v1/usage/report`<br>
Reports usage event (called by your main API service)<br>
Headers: `Authorization: Bearer <api_key>`<br>
Body: 
```js 
{ 
  event_type: "api_call",
  quantity: 1, 
  endpoint: "/v1/generate", 
  metadata: { /* ... */} 
}
```
Response: 
```js
{ 
  success: true, 
  data: {
    recorded_at: "...",
    remaining_quota: {
      tokens: 35000, 
      requests: 450
    } 
  } 
}
```

`POST /api/v1/usage/batch-report`<br>
Reports multiple usage events in a single call (for efficiency)<br>
Headers: `Authorization: Bearer <api_key>`<br>
Body: 
```js
{ 
  events: [{
    event_type: "tokens_used", 
    quantity: 1500,
    timestamp: "..." 
  },
  { /* ... */ }] 
}
```
Response: 
```js
{ success: true, data: { processed_count: 10, failed_count: 0 } }
```

`GET /api/v1/usage/quota-check`<br>
Checks if user has remaining quota before processing request<br>
Headers: `Authorization: Bearer <api_key>`<br>
Query params: `?tokens_needed=1000&requests_needed=1`<br>
Response: 
```js
{ 
  success: true,
  data: {
    allowed: true,
    remaining: {
      tokens: 35000,
      requests: 450
    }
  }
}
```
Response (over quota): 
```js
{ 
  success: false,
  data: "Quota exceeded", 
  quota_info: { /* ... */ } 
}
```


## Webhook Endpoints (No Authentication, Signature Verified)

#### Payment Gateway Webhooks
`POST /webhooks/stripe`<br>
Handles Stripe webhook events (subscription updates, payment failures, etc.)<br>
Headers: `Stripe-Signature: <webhook_signature>`<br>
Body: Stripe webhook payload<br>
Events handled: `customer.subscription.created, invoice.payment_succeeded, invoice.payment_failed, etc.`<br>
Response: 
```js
{ "received": true }
```

`POST /webhooks/xendit`<br>
Handles Xendit webhook events (for Indonesian market)<br>
Headers: `x-callback-token: <xendit_callback_token>`<br>
Body: Xendit webhook payload<br>
Response: 
```js
{ "received": true }
```

`POST /webhooks/midtrans`<br>
Handles Midtrans webhook events (for Indonesian market)<br>
Headers: `Authorization: Basic <base64(server_key:)>`<br>
Body: Midtrans webhook payload<br>
Response:
```js
{ "received": true }
```

## Administrative Endpoints (Admin Authentication Required)

#### Admin User Management
`GET /admin/customers`<br>
Lists all customers with pagination and filters<br>
Headers: `Authorization: Bearer <admin_jwt>`<br>
Query params: `?search=john@example.com&status=active&limit=50&offset=0`<br>
Response: 
```js
{ 
  success: true,
    "data": { "customers": [...], "total": 1234, "page": 1 } }
```

`GET /admin/customers/{id}`<br>
Gets detailed customer information including usage, subscriptions, etc.<br>
Headers: `Authorization: Bearer <admin_jwt>`<br>
Response: 
```js
{ 
  success: true,
  data: { 
    customer: { /* ... */},
    subscription: {/* ... */},
    usage_stats: {/* ... */}
  } 
}
```

`PUT /admin/customers/{id}/subscription`<br>
Manually update customer subscription (admin override)<br>
Headers: `Authorization: Bearer <admin_jwt>`<br>
Body: 
```js
{ plan_id: "uuid", status: "active", extend_trial: 7 }
```
Response:
```js
{ success: true, data: { subscription: { /* ... */} } }
```

`POST /admin/customers/{id}/credit`<br>
Add billing credit to customer account<br>
Headers: `Authorization: Bearer <admin_jwt>`<br>
Body: 
```js
{ amount_cents: 5000, reason: "Service disruption compensation" }
```
Response: 
```js
{ success: true, data: { credit_applied: 5000, new_balance: 5000 } }
```


#### Admin Analytics & Reports
`GET /admin/analytics/overview`<br>
Gets high-level business metrics<br>
Headers: `Authorization: Bearer <admin_jwt>`<br>
Query params: `?period=last_30_days`<br>
Response: 
```js
{ 
  success: true,
  data: { 
    total_revenue: 150000,
    active_subscriptions: 45,
    churn_rate: 0.05
  }
}
```

`GET /admin/analytics/usage`<br>
Gets aggregated usage statistics across all customers<br>
Headers: `Authorization: Bearer <admin_jwt>`<br>
Response: 
```js
{ 
  success: true,
  data: {
    total_requests: 1000000,
    avg_requests_per_customer: 2222
  }
}
```

`GET /admin/invoices/failed`<br>
Lists failed payment invoices requiring attention<br>
Headers: `Authorization: Bearer <admin_jwt>`<br>
Response: 
```js
{ 
  success: true,
  data: [{
    invoice_id: "uuid",
    customer: {/* ... */},
    amount: 2999,
    failure_reason: "..." 
  }]
}
```

#### Admin Configuration
`GET /admin/plans`<br>
Lists all plans including inactive ones<br>
Headers: `Authorization: Bearer <admin_jwt>`<br>
Response: 
```js
{ 
  success: true,
  data: [{
    id: "uuid",
    name: "Pro",
    is_active: false,
    //... 
  }] 
}
```

`POST /admin/plans`<br>
Creates a new subscription plan<br>
Headers: `Authorization: Bearer <admin_jwt>`<br>
Body: 
```js
{ 
  name: "Enterprise", 
  price_cents: 9999,
  features: { /* ... */ },
  quota_tokens: 100000 
}
```
Response: 
```js
{ 
  success: true,
  data: {
    id: "uuid",
    name: "Enterprise", 
    //... 
  }
}
```

`PUT /admin/plans/{id}`<br>
Updates existing plan (careful with active subscriptions)<br>
Headers: `Authorization: Bearer <admin_jwt>`<br>
Body: 
```js
{ price_cents: 3499, quota_tokens: 75000 }
```
Response: 
```js
{ 
  success: true,
  data: {
    plan: {/* ... */},
    affected_subscriptions: 12
  } 
}
```


### Health & Monitoring
`GET /health`<br>
Basic health check endpoint<br>
Response: 
```js
{ success: true, data: "OK" }
```

`GET /health/detailed`<br>
Detailed health check including database, external services<br>
Response: 
```js
{ 
  success: true,
  data: { 
    database: "OK",
    stripe: "OK",
    redis: "OK",
    uptime: "72h"
  } 
}
```

`GET /metrics`<br>
Prometheus metrics endpoint (if implemented)<br>
Response: Prometheus format metrics
