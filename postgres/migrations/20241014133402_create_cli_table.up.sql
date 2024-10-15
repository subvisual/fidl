BEGIN;

CREATE TABLE IF NOT EXISTS
  cli (
    id bigint GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    signature text NOT NULL,
    balance float NOT NULL DEFAULT 0.0,
    created_at timestamp(0) NOT NULL DEFAULT (NOW() at time zone ('utc')),
    updated_at timestamp(0) NOT NULL DEFAULT (NOW() at time zone ('utc'))
  );

CREATE UNIQUE INDEX cli_idx ON cli (signature);

COMMIT;