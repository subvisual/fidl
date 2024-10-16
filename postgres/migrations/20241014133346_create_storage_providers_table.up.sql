BEGIN;

CREATE TABLE IF NOT EXISTS
  storage_provider (
    id bigint GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    sp_id text NOT NULL,
    wallet_address text NOT NULL,
    price float NOT NULL DEFAULT 0.0,
    created_at timestamp(0) NOT NULL DEFAULT (NOW() at time zone ('utc')),
    updated_at timestamp(0) NOT NULL DEFAULT (NOW() at time zone ('utc'))
  );

CREATE UNIQUE INDEX storage_provider_idx ON storage_provider (sp_id, wallet_address);

COMMIT;