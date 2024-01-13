-- name: FindUser :one
SELECT * FROM "users" WHERE id=$1;

-- name: CheckWalletExists :one
SELECT EXISTS (
  SELECT id FROM "wallet" WHERE user_id=$1
);

-- name: FindTokenByUserId :one
SELECT * FROM "users" WHERE token=$1;

-- name: CreateWallet :one
INSERT INTO "wallet"(id, user_id, status, balance)
VALUES($1, $2, $3, $4) RETURNING *;

-- name: FindWalletByUserId :one
SELECT * FROM "wallet" WHERE user_id=$1;

-- name: UpdateWalletByUserId :one
UPDATE "wallet" 
SET 
  status=$1,
  enabled_at=NOW()
WHERE user_id=$2 RETURNING *;

-- name: FindTransactionByUserId :many
SELECT * FROM "transaction" WHERE user_id=$1 ORDER BY created_at ASC;

-- name: UpdatBalanceWalletByUserId :one
UPDATE "wallet" 
SET 
  balance=$1
WHERE user_id=$2 RETURNING *;

-- name: CreateTransaction :one
INSERT INTO "transaction"(id, user_id, status, type, balance, reference_id)
VALUES($1, $2, $3, $4, $5, $6) RETURNING *;
