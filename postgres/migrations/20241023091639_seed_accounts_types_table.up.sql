BEGIN;

INSERT INTO
  accounts_types (id, name)
VALUES
  (1, 'Storage Provider'),
  (2, 'Client');

COMMIT;
