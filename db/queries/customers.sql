-- name: CreateCustomer :one
INSERT INTO customers (id, user_id, email, credit_balance_cents)
VALUES ($1, $2, $3, 0)
RETURNING *;

-- name: GetCustomerByUserID :one
SELECT * FROM customers WHERE user_id = $1 LIMIT 1;

-- name: UpdateCustomerCredits :one
UPDATE customers
SET credit_balance_cents = credit_balance_cents + $2
WHERE id = $1
RETURNING *;
