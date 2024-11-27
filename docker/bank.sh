#!/bin/bash

set -eux

PUID=${PUID:-1000}
PGID=${PGID:-1000}

groupmod -o -g "$PGID" runner &> /dev/null
usermod -o -u "$PUID" runner &> /dev/null

chown -R runner:runner /app

cmd="$*"
if [ "$cmd" = "/bin/bash" ]; then
  exec gosu runner "$@"
  exit
fi

DSN=$(sed -nr 's/dsn="(.*)"/\1/p' /etc/bank.ini | tr -d [:space:])
/app/migrate -path=/app/postgres/migrations -database=$DSN up
exec gosu runner "$@"
