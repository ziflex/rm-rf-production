# Accounts & Transactions API (Go + Echo)

A small REST API that manages customer accounts and transactions with strict validation and OpenAPI‑first development.

- Language: Go
- Framework: Echo
- API contract and validation: OpenAPI 3 + oapi‑codegen
- Database: PostgreSQL
- Containerization: Docker Compose
- Docs: Swagger UI

## Quick start

### 1) Run with Docker Compose

```bash
make up
```

App will listen on http://localhost:8080

Services:
- `db` – PostgreSQL 15
- `app` – Go API server
- `migrations` – runs SQL migrations at startup

To shut down:

```bash
make down
```

#### How to run locally without Docker
1. Install dependencies: `make install`
2. Build the app: `make build`
3. Start a local Postgres instance (e.g. `docker run -p 5432:5432 -e POSTGRES_USER=user -e POSTGRES_PASSWORD=password -e POSTGRES_DB=mydb -d postgres:15`)
4. Run migrations: `DB_HOST=localhost DB_PORT=5432 DB_USER=user DB_PASS=password DB_NAME=mydb make migrate`
5. Start the app: `DB_HOST=localhost DB_PORT=5432 DB_USER=user DB_PASS=password DB_NAME=mydb make start`

### 2) Explore the API
- Swagger UI: `http://localhost:8080/docs`
- OpenAPI spec: `http://localhost:8080/openapi.yaml`

### 3) Health check

```bash
curl -i http://localhost:8080/health
```

## Configuration

The server uses environment variables (parsed via `caarlos0/env`). These are already provided by `docker-compose.yaml`, but can be overridden.

| Var         | Default     | Description                 |
|-------------|-------------|-----------------------------|
| `PORT`      | `8080`      | HTTP port                   |
| `LOG_LEVEL` | `trace`     | zerolog level               |
| `DB_HOST`   | `localhost` | Postgres host               |
| `DB_PORT`   | `5432`      | Postgres port               |
| `DB_NAME`   | `app`       | Database name               |
| `DB_USER`   | `app`       | Database user               |
| `DB_PASS`   | `app`       | Database password           |

Example Compose service block for the app:
```yaml
environment:
  PORT: "8080"
  LOG_LEVEL: trace
  DB_HOST: db
  DB_PORT: "5432"
  DB_NAME: app
  DB_USER: app
  DB_PASS: app
```

## API overview

### Create account
Create a new customer account by unique `document_number`.

```
POST /accounts
Content-Type: application/json
```

Request
```json
{
  "document_number": "12345678900"
}
```

201 Created
```json
{
  "id": 1,
  "document_number": "12345678900"
}
```

Errors
- 400 invalid payload
- 409 document number already exists

---

### Get account by id
Fetch an existing account.

```
GET /accounts/{accountId}
```

200 OK
```json
{
  "id": 1,
  "document_number": "12345678900"
}
```

Errors
- 404 account not found

---

### Create transaction
Create a transaction for an account. Client must send a **positive** `amount`; the server applies the proper sign when storing.

```
POST /transactions
Content-Type: application/json
```

Operation types:
- `PURCHASE`
- `INSTALLMENT`
- `WITHDRAWAL`
- `PAYMENT`

Request
```json
{
  "account_id": 1,
  "operation_type": "PURCHASE",
  "amount": 100.00
}
```

201 Created
```json
{
  "id": 1,
  "account_id": 1,
  "operation_type": "PURCHASE",
  "amount": -100.00,
  "event_date": "2025-08-30T19:49:41Z"
}
```

Errors
- 400 invalid payload or operation type
- 404 account not found

---

## cURL examples

Create account
```bash
curl -sS -X POST http://localhost:8080/accounts   -H 'Content-Type: application/json'   -d '{"document_number":"12345678900"}' | jq
```

Get account
```bash
curl -sS http://localhost:8080/accounts/1 | jq
```

Create transaction (client sends positive amount)
```bash
curl -sS -X POST http://localhost:8080/transactions   -H 'Content-Type: application/json'   -d '{"account_id":1,"operation_type":"PURCHASE","amount":100.00}' | jq
```

---

## Project layout

```
.
├── database/
│   ├── migrations/         # Database migration files
│   ├── Dockerfile          # Dockerfile for migration service
│   └── entrypoint.sh       # Migration runner
├── internal/
│   ├── api/                # HTTP handlers, routing, middleware
│   ├── database/           # DB wiring, repositories
│   └── server/             # Echo server bootstrap
├── pkg/
│   ├── accounts/           # Domain model + service
│   └── transactions/       # Domain model + service
├── spec/
│   ├── ui/                 # Swagger UI assets (served at /docs)
│   ├── file.go             # Embedded OpenAPI spec and UI assets
│   ├── oapi-codegen.go     # Configuration for oapi-codegen
│   └── openapi.yaml        # API contract (served at /openapi.yaml)
└── main.go                 # App entrypoint
```

## Database schema

Tables
- `accounts(id serial primary key, document_number text unique not null)`
- `transactions(id serial primary key, account_id int not null references accounts(id), operation_type enum not null, amount numeric not null, event_date timestamp not null default now())`

Indexes
- `transactions(account_id)`

Enum
- `operation_type` with the 4 values listed above.


## Development

### Prerequisites
- Go 1.23+ (or the version pinned in `go.mod`)
- Docker (for local DB)
- Make (optional)

### Run locally
```bash
make install
make build
make start
```

### Tests
```bash
make test
```

## Design notes
- OpenAPI‑first: the server mounts the spec at `/openapi.yaml` and uses validation middleware to enforce request and response shapes.
- Amount sign is applied on the server for consistency and client simplicity.
- Errors use consistent JSON payloads and proper HTTP status codes.
- Timestamps are in UTC ISO‑8601.