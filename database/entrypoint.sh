#!/bin/sh

set -e

: "${DB_HOST:?Need to set DB_HOST}"
: "${DB_PORT:?Need to set DB_PORT}"
: "${DB_USER:?Need to set DB_USER}"
: "${DB_PASSWORD:?Need to set DB_PASSWORD}"
: "${DB_NAME:?Need to set DB_NAME}"

DATABASE_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"

exec migrate -path ./migrations -database "$DATABASE_URL" up