-- name: CreatePlan :one
INSERT INTO plans (id, slug, name, description, price_cents, currency, interval, quota_limits, meta)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: ListPlans :many
SELECT * FROM plans;

-- name: GetPlanBySlug :one
SELECT * FROM plans
WHERE slug = $1
LIMIT 1;

