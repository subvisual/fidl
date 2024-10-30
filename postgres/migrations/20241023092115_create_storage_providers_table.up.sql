BEGIN;

CREATE TABLE IF NOT EXISTS
  storage_providers (
    id bigint PRIMARY KEY REFERENCES accounts(id),
    sp_id text NOT NULL,
    price numeric(38) NOT NULL DEFAULT 0,
    created_at timestamp(0) NOT NULL DEFAULT (NOW() at time zone ('utc')),
    updated_at timestamp(0) NOT NULL DEFAULT (NOW() at time zone ('utc'))
  );

COMMIT;