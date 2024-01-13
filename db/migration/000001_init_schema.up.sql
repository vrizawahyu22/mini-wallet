CREATE TABLE IF NOT EXISTS "users" (
  "id" varchar PRIMARY KEY,
  "token" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE IF NOT EXISTS "wallet" (
  "id" varchar PRIMARY KEY,
  "user_id" varchar NOT NULL,
  "status" varchar NOT NULL DEFAULT 'disabled',
  "balance" numeric NOT NULL,
  "enabled_at" timestamptz NOT NULL DEFAULT (now()),
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE IF NOT EXISTS "transaction" (
  "id" varchar PRIMARY KEY,
  "user_id" varchar NOT NULL,
  "status" varchar NOT NULL,
  "type" varchar NOT NULL,
  "balance" numeric NOT NULL,
  "reference_id" varchar UNIQUE NOT NULL,
  "transaction_at" timestamptz NOT NULL DEFAULT (now()),
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "users" ("token");
CREATE INDEX ON "wallet" ("user_id");

ALTER TABLE "wallet" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;
ALTER TABLE "transaction" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;

INSERT INTO "users"("id", "token")
VALUES('ea0212d3-abd6-406f-8c67-868e814a2436', '6b3f7dc70abe8aed3e56658b86fa508b472bf238');
