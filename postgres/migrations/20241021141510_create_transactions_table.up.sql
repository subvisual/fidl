BEGIN;

CREATE TABLE IF NOT EXISTS
  transactions (
    id bigint GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    uuid UUID DEFAULT gen_random_uuid (),
    source text NOT NULL,
    destination text NOT NULL,
    value float NOT NULL,
    status_id integer NOT NULL DEFAULT 1 REFERENCES transaction_status (id),
    created_at timestamp(0) NOT NULL DEFAULT (NOW() at time zone ('utc')),
    updated_at timestamp(0) NOT NULL DEFAULT (NOW() at time zone ('utc'))
  );

CREATE INDEX transaction_status_idx ON transactions (status_id);

COMMIT;