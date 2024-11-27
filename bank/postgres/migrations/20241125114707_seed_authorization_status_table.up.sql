BEGIN;

INSERT INTO
  authorization_status (id, name)
VALUES
  (1, 'Open'),
  (2, 'Locked');

COMMIT;