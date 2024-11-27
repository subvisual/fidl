BEGIN;

CREATE TABLE IF NOT EXISTS
  escrow (
    id bigint REFERENCES accounts(id),
    uuid UUID PRIMARY KEY NOT NULL,
    balance numeric(38) NOT NULL DEFAULT 0,
    proxy text NOT NULL,
    status_id integer NOT NULL DEFAULT 1 REFERENCES authorization_status (id),
    created_at timestamp(0) NOT NULL DEFAULT (NOW() at time zone ('utc')),
    updated_at timestamp(0) NOT NULL DEFAULT (NOW() at time zone ('utc'))
  );

COMMIT;