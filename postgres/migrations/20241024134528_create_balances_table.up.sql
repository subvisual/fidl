BEGIN;

CREATE TABLE IF NOT EXISTS
  balances (
    id bigint PRIMARY KEY REFERENCES accounts(id),
    balance numeric(30,18) NOT NULL DEFAULT 0.0,
    created_at timestamp(0) NOT NULL DEFAULT (NOW() at time zone ('utc')),
    updated_at timestamp(0) NOT NULL DEFAULT (NOW() at time zone ('utc'))
  );

COMMIT;