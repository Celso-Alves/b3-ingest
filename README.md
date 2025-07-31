# b3-ingest

It's a simple Go project for automated ingestion and querying of B3 trading data

## Project Overview

**b3-ingest** is an ETL/data pipeline for ingesting, processing, and serving B3 (Brazilian Stock Exchange) trading data. It is designed for maintainability, performance, and extensibility, using:

- **Clean Architecture**: Clear separation of domain, service, repository, and infrastructure layers.
- **SOLID Principles**: Each component has a single responsibility and is easily testable and replaceable.
- **Idiomatic Go**: Uses context, dependency injection, and best practices for concurrency and error handling.
- **Environment-driven configuration**: All settings are loaded from environment variables or `.env` files.
- **Dockerized**: For reproducible local development and deployment.

## Architecture

```
main.go
  └── internal/
      ├── domain/         # Domain models (pure Go, no dependencies)
      ├── service/        # Business logic (ingestion, trading, etc.)
      ├── infra/          # Infrastructure (DB, repositories, adapters)
      ├── logger/         # Custom logger abstraction
      ├── starter/        # Application startup orchestration
      └── settings/       # Environment/config loading
  └── pkg/routes/         # HTTP route handlers (Gin)
  └── cmd/                # Compiled binary output
```

- **Domain and ORM models are decoupled** for testability and future DB changes.
- **Repository pattern**: All DB access is via repository interfaces, using GORM.
- **Service layer**: All business logic (ingestion, trading queries) is in services, not handlers or main.


## Prerequisites

- Go 1.24 or newer (recommended: latest stable)
- Docker & Docker Compose (for local DB and app)


## Environment Setup

1. **Clone the repository:**
   ```sh
   git clone https://github.com/Celso-Alves/b3-ingest.git
   cd b3-ingest
   ```
2. **Dependencies:**


The first step to running the project will be to configure the environment, installing the necessary dependencies, this can be done through the command in your terminal:

```bash
go mod download
```

3. **Configure environment variables(optional):**
   - Copy and edit `env-local` as needed:
     ```sh
     cp env-local .env
     # Edit .env to match your local DB settings if needed
     ```
   - Example `.env`:
     ```env
     DATABASE_NAME=b3db
     DATABASE_HOST=localhost
     DATABASE_PORT=5432
     DATABASE_USERNAME=postgres
     DATABASE_PASSWORD=postgres
     APP_NAME=b3-ingest
     ```
    - Load the envs before running:
      ```bash
      source setenv.sh
      ```


### Build DB dependencies with Docker Compose (app + Postgres)

```sh
make docker-up
```
- You can access the API at [http://localhost:8000](http://localhost:8000)


### Build the application

```sh
make build
```
- The binary will be generated in `cmd/b3-ingest`.

### Download the latest B3 trading files (last 7 workdays)

```sh
make download
```
- Files will be downloaded and unzipped to `bundle/b3files` (default).

### Run the ingestion (load CSVs into the database)

```sh
make ingest
```

### Run the HTTP server

```sh
make serve
```
- The server will start on the port defined in your environment (default: 8000).


## Querying the Data

After running the server (`make serve` or via Docker), you can query for trading stats using HTTP:

### Example: Query max price and max daily volume for a ticker

```sh
curl "http://localhost:8000/quote?ticker=WDOQ25&data_inicio=2025-07-29"
```
Response:
```json
{
  "ticker": "WDOQ25",
  "max_range_value": 5585.0,
  "max_daily_volume": 4688104
}
```
- `ticker` (required): The instrument code.
- `data_inicio` (optional, YYYY-MM-DD): Start date for the query (default: 7 days ago).

## Best Practices Used

- **Clean Architecture**: All business logic is in services, not in handlers or main.
- **SOLID Principles**: Each layer/component has a single responsibility and is easily testable.
- **Idiomatic Go**: Uses context, error wrapping, dependency injection, and concurrency best practices.
- **Graceful shutdown**: All modes handle OS signals and shutdown cleanly.
- **Environment-driven config**: All config is loaded from env or `.env` files, never hardcoded.
- **Automated tests**: Unit tests for all core logic, with Arrange/Act/Assert and GivenWhenThen naming.


## Configuration

The application is configured via environment variables. Below is a summary of all supported variables:

| Name                | Description                                 | Default Value           |
|---------------------|---------------------------------------------|------------------------|
| `CSV_PATH`          | Directory where CSV files are stored        | `./bundle/b3files`     |
| `APP_DEFAULT_PORT`  | HTTP server port                            | `8000`                 |
| `APP_NAME`          | Application name                            | `b3-ingest`            |
| `INGESTION_CORES`   | Number of concurrent ingestion workers      | `6`                    |
| `DATABASE_NAME`     | PostgreSQL database name                    | `b3db`                 |
| `DATABASE_PASSWORD` | PostgreSQL user password                    | `postgres`             |
| `DATABASE_USERNAME` | PostgreSQL username                         | `postgres`             |
| `DATABASE_HOST`     | PostgreSQL host                             | `localhost`            |
| `DATABASE_PORT`     | PostgreSQL port                             | `5432`                 |
| `DATABASE_SSL`      | Use SSL for DB connection (`true`/`false`)  | `false`                 |

- All variables can be set in your shell, `.env`, or via Docker Compose.
- Required variables must be set for the application to start.

