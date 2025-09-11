<div align="center">
  <h1>⚡ TADB API</h1>
  <p>
    <img src="https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go Version">
    <img src="https://img.shields.io/badge/Gin-1.10.1-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Gin Framework">
    <img src="https://img.shields.io/badge/PostgreSQL-15+-336791?style=for-the-badge&logo=postgresql&logoColor=white" alt="PostgreSQL">
    <img src="https://img.shields.io/badge/License-MIT-green?style=for-the-badge" alt="License">
  </p>

  <p><em>A RESTful API for managing energy matrix data with generators, types, and production tracking.</em></p>
</div>

# Requirements

- `Go >= 1.21`
- `PostgreSQL >= 15`

# Database Schema

The TADB API uses a PostgreSQL database with the following structure:

## Core Schema Tables

### `core.type`
Stores different types of energy generators (renewable/non-renewable).
```sql
- id (UUID, Primary Key)
- name (VARCHAR(20), Unique) - Type name
- description (VARCHAR(80)) - Type description  
- isRenuevable (BOOLEAN) - Whether the type is renewable
```

### `core.generator`
Stores individual energy generators with their capacity.
```sql
- id (UUID, Primary Key)
- type (UUID, Foreign Key → core.type.id)
- capacity (FLOAT) - Generator capacity in MW
```

### `core.production` 
Tracks daily energy production for each generator.
```sql
- id (UUID, Primary Key)
- generator_id (UUID, Foreign Key → core.generator.id)
- date (DATE) - Production date
- production_mw (DECIMAL) - Production in megawatts
- UNIQUE(generator_id, date) - One record per generator per day
```

## Relationships
- **Type → Generator**: One-to-Many (one type can have multiple generators)
- **Generator → Production**: One-to-Many (one generator can have multiple production records)

# Installation

1. Clone the repository
```bash
git clone https://github.com/02loveslollipop/api_matriz_enegertica_tadb.git
cd api_tadb
```

2. Set up PostgreSQL database
```bash
# Create database
createdb tadb

# Run the schema creation script
psql -d tadb -f sql/create.sql
```

3. Install Go dependencies
```bash
go mod download
go mod tidy
```

4. Set up environment variables
```bash
# Create .env file
cp .env.example .env

# Edit with your database configuration
DB_HOST=localhost
DB_PORT=5432
DB_NAME=tadb
DB_USER=your_username
DB_PASSWORD=your_password
```

5. Run the application
```bash
# Using Go directly
go run cmd/main.go

# Or using Make
make run

# Or build and run
make build
./bin/tadb-api
```

The server will start on `http://localhost:8080`

# Features

## API Endpoints

### Core Endpoints
- `GET /` - Welcome message and API info
- `GET /health` - Health check endpoint

### Generator Types
- `GET /api/v1/types` - List all generator types
- `GET /api/v1/types/:id` - Get specific type
- `POST /api/v1/types` - Create new type
- `PUT /api/v1/types/:id` - Update type
- `DELETE /api/v1/types/:id` - Delete type

### Generators
- `GET /api/v1/generators` - List all generators
- `GET /api/v1/generators/:id` - Get specific generator
- `POST /api/v1/generators` - Create new generator
- `PUT /api/v1/generators/:id` - Update generator
- `DELETE /api/v1/generators/:id` - Delete generator

### Production Data
- `GET /api/v1/productions` - List production records
- `GET /api/v1/productions/:id` - Get specific production record
- `POST /api/v1/productions` - Create production record
- `PUT /api/v1/productions/:id` - Update production record
- `DELETE /api/v1/productions/:id` - Delete production record

### Analytics Endpoints
- `GET /api/v1/analytics/total-production` - Total production by date range
- `GET /api/v1/analytics/renewable-vs-nonrenewable` - Renewable vs non-renewable production
- `GET /api/v1/analytics/generator-efficiency` - Generator efficiency metrics

# License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

# Acknowledgments

- [Gin Web Framework](https://github.com/gin-gonic/gin) - HTTP web framework
- [PostgreSQL](https://www.postgresql.org/) - Database system
- [UUID Extension](https://www.postgresql.org/docs/current/uuid-ossp.html) - UUID generation

## API Documentation

The project generates OpenAPI docs from code annotations using `swaggo/swag` and serves Swagger UI.

- Local Swagger UI: `http://localhost:8080/swagger/index.html`
- Generate spec locally:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
swag init -g cmd/main.go -d cmd,pkg -o docs
```

The generated files are written to `docs/swagger.json` and `docs/swagger.yaml`. The CI workflow publishes `docs/swagger.yaml` to Bump.sh.
