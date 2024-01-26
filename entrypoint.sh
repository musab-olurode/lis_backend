#!/bin/sh

# wait for Postgres to start
function postgres_ready() {
python << END
import sys
import psycopg2
try:
    dbUrl = "$DB_URL"
    conn = psycopg2.connect(dbUrl)
except psycopg2.OperationalError:
    sys.exit(-1)
sys.exit(0)
END
}

until postgres_ready; do
  >&2 echo "Postgres is unavailable - sleeping"
  sleep 1
done

# run migrations
make migrate

# run the main application
exec "$@"