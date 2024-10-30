BEGIN;

INSERT INTO
  balance_status (id, name)
VALUES
  (1, 'Available'),
  (2, 'Locked');

COMMIT;
