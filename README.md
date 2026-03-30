# Power-Pi

A structured REST API written in Go using Gorilla Mux, GORM, and PostgreSQL, fully containerised with Docker.

## Stack

| Layer     | Technology                                    |
| --------- | --------------------------------------------- |
| Router    | [gorilla/mux](https://github.com/gorilla/mux) |
| ORM       | [GORM](https://gorm.io) + PostgreSQL driver   |
| Database  | PostgreSQL 17                                 |
| Logging   | [zerolog](https://github.com/rs/zerolog)      |
| Config    | [godotenv](https://github.com/joho/godotenv)  |
| Container | Docker + Docker Compose                       |

## Project Structure

```
Power-Pi/
├── main.go                 # Entry point
├── .env                    # Local environment variables
├── Dockerfile              # Multi-stage build
├── docker-compose.yml      # App + database services
├── config/
│   └── config.go           # Loads env variables with defaults
├── database/
│   ├── db.go               # Postgres connection with connection pool + AutoMigrate
│   └── models.go           # GORM models (PowerTable)
├── apis/
│   ├── routes.go           # Gorilla Mux router setup
│   ├── handlers.go         # HTTP handlers
│   └── middleware.go       # Logging middleware
└── logger/
    └── logger.go           # Zerolog initialisation
```

## Environment Variables

Configure the application via `.env`:

| Variable       | Default     | Description                                                    |
| -------------- | ----------- | -------------------------------------------------------------- |
| `DB_HOST`      | `localhost` | Database host (use `db` inside Docker)                         |
| `DB_PORT`      | `5432`      | Database port (internal container port)                        |
| `DB_HOST_PORT` | `5432`      | Host port mapped to the database                               |
| `DB_USER`      | `postgres`  | Database user                                                  |
| `DB_PASSWORD`  | `postgres`  | Database password                                              |
| `DB_NAME`      | `powerpi`   | Database name                                                  |
| `SERVER_PORT`  | `8080`      | API server port                                                |
| `LOG_LEVEL`    | `info`      | Log level (`trace`, `debug`, `info`, `warn`, `error`, `fatal`) |

## Getting Started

### Run with Docker (recommended)

```bash
docker compose up --build -d
```

Both the app and database will start. The API will be available at `http://localhost:9090`.

To stop and remove containers and volumes:

```bash
docker compose down -v
```

### Run locally

Requires a running PostgreSQL instance. Update `.env` with `DB_HOST=localhost`, then:

```bash
go run main.go
```

## API Endpoints

### `GET /power-table`

Returns all records.

```bash
curl http://localhost:9090/power-table
```

**Response `200 OK`**
```json
[
  {
    "ID": 1,
    "CreatedAt": "...",
    "UpdatedAt": "...",
    "DeletedAt": null,
    "price": 99.99,
    "company": "Acme Corp"
  }
]
```

---

### `POST /power-table`

Creates a new record.

```bash
curl -X POST http://localhost:9090/power-table \
  -H "Content-Type: application/json" \
  -d '{"price": 99.99, "company": "Acme Corp"}'
```

**Response `201 Created`**
```json
{
  "ID": 1,
  "CreatedAt": "...",
  "UpdatedAt": "...",
  "DeletedAt": null,
  "price": 99.99,
  "company": "Acme Corp"
}
```

## Connection Pool

Configured in `database/db.go`:

| Setting              | Value     |
| -------------------- | --------- |
| `SetMaxIdleConns`    | 10        |
| `SetMaxOpenConns`    | 100       |
| `SetConnMaxLifetime` | 5 minutes |

## Database Migrations

GORM `AutoMigrate` runs automatically on startup. It creates or updates the `power_tables` table to match the current model without dropping existing data.
