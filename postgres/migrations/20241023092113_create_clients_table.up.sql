BEGIN;

CREATE TABLE IF NOT EXISTS
  clients (
    id bigint PRIMARY KEY REFERENCES accounts(id),
    status_id integer NOT NULL DEFAULT 1 REFERENCES balance_status (id),
    created_at timestamp(0) NOT NULL DEFAULT (NOW() at time zone ('utc')),
    updated_at timestamp(0) NOT NULL DEFAULT (NOW() at time zone ('utc'))
  );

CREATE INDEX balance_status_idx ON clients (status_id);

COMMIT;