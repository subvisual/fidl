#!/bin/bash

set -eux

PUID=${PUID:-1000}
PGID=${PGID:-1000}

groupmod -o -g "$PGID" runner &> /dev/null
usermod -o -u "$PUID" runner &> /dev/null

chown -R runner:runner /app

exec gosu runner "$@"
