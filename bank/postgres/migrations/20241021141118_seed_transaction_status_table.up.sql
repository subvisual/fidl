BEGIN;

INSERT INTO
  transaction_status (id, name)
VALUES
  (1, 'Pending'),
  (2, 'Completed');

COMMIT;
