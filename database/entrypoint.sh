#!/bin/sh

set -e

: "${DB_HOST:?Need to set DB_HOST}"
: "${DB_PORT:?Need to set DB_PORT}"
: "${DB_USER:?Need to set DB_USER}"
: "${DB_PASS:?Need to set DB_PASS}"
: "${DB_NAME:?Need to set DB_NAME}"

DATABASE_URL="postgres://${DB_USER}:${DB_PASS}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"

until nc -z "$DB_HOST" "$DB_PORT"; do
  echo "Waiting for Postgres at $DB_HOST:$DB_PORT..."
  sleep 1
done

exec migrate -path ./migrations -database "$DATABASE_URL" up