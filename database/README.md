# CMDB Lite Database

This directory contains the database schema, migrations, and seed data for CMDB Lite.

## Database Schema

The database consists of the following tables:

- `users` - User accounts for authentication
- `configuration_items` - The actual configuration items
- `relationships` - Relationships between configuration items
- `audit_logs` - Audit logs for tracking changes

## Getting Started

### Prerequisites

- PostgreSQL 12 or higher
- psql command-line tools

### Installation

1. Copy the environment file:
   ```bash
   cp .env.example .env
   ```

2. Create the database:
   ```bash
   ./migrate.sh create
   ```

3. Apply the schema:
   ```bash
   ./migrate.sh schema
   ```

4. Apply migrations:
   ```bash
   ./migrate.sh migrate
   ```

5. Apply seed data:
   ```bash
   ./migrate.sh seed
   ```

### Migration Script

The `migrate.sh` script provides the following commands:

- `create` - Create the database
- `schema` - Apply database schema
- `migrate` - Apply database schema and migrations
- `seed` - Apply database schema, migrations, and seed data

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| DB_HOST | Database host | localhost |
| DB_PORT | Database port | 5432 |
| DB_NAME | Database name | cmdb_lite |
| DB_USER | Database username | cmdb_user |
| DB_PASSWORD | Database password | cmdb_password |

## Project Structure

```
database/
├── migrations/            # Database migration files
│   └── 001_add_initial_data.sql
├── schema/                # Database schema files
│   └── 001_initial_schema.sql
├── seeds/                 # Database seed files
│   └── 001_initial_data.sql
├── .env.example           # Environment variables example
└── migrate.sh             # Migration script
```

## Schema Details

### Users Table

| Column | Type | Description |
|--------|------|-------------|
| id | UUID | Primary key |
| username | VARCHAR(50) | Unique username |
| password_hash | VARCHAR(255) | Hashed password |
| role | VARCHAR(20) | User role (default: 'viewer') |
| created_at | TIMESTAMP | Creation timestamp |

### Configuration Items Table

| Column | Type | Description |
|--------|------|-------------|
| id | UUID | Primary key |
| name | VARCHAR(255) | CI name |
| type | VARCHAR(100) | CI type |
| attributes | JSONB | Flexible attributes as JSON |
| tags | TEXT[] | Array of tags |
| created_at | TIMESTAMP | Creation timestamp |
| updated_at | TIMESTAMP | Last update timestamp |

### Relationships Table

| Column | Type | Description |
|--------|------|-------------|
| id | UUID | Primary key |
| source_id | UUID | Reference to source CI |
| target_id | UUID | Reference to target CI |
| type | VARCHAR(100) | Relationship type |
| created_at | TIMESTAMP | Creation timestamp |

### Audit Logs Table

| Column | Type | Description |
|--------|------|-------------|
| id | UUID | Primary key |
| entity_type | VARCHAR(50) | Type of entity (e.g., 'configuration_item') |
| entity_id | UUID | ID of the entity |
| action | VARCHAR(50) | Action performed (e.g., 'create', 'update', 'delete') |
| changed_by | VARCHAR(50) | Username of who made the change |
| changed_at | TIMESTAMP | When the change was made |
| details | JSONB | Details of the change |

## Goose Migrations

The backend also uses Goose for database migrations. The migration files are located in `backend/migrations/` and follow the naming convention `YYYYMMDDHHMMSS_description.up.sql` and `YYYYMMDDHHMMSS_description.down.sql`.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Commit your changes
6. Push to the branch
7. Create a Pull Request