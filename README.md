# dom-cobb

Dom Cobb dives deep into the layered world of feature flags. This backend service manages feature toggles with support for dependencies, ensures all prerequisites are met before activation, and logs every action with full audit trails — all while keeping your system from slipping into recursive nightmares like circular dependencies.

## Features

- **Feature Flag Management**: Create, toggle, and manage feature flags with comprehensive validation
- **Dependency Support**: Define hierarchical dependencies between flags with circular dependency detection
- **Audit Logging**: Complete audit trail of all operations with timestamps, reasons, and actor information
- **Validation Engine**: Prevents invalid state changes by validating dependencies before flag operations
- **RESTful API**: Clean, well-documented API endpoints for all operations
- **Dockerized**: Fully containerized with Docker Compose for easy deployment
- **Testing Suite**: Comprehensive test coverage with Ginkgo testing framework

## Tech Stack

- **Backend**: Go with Gin framework
- **Primary Database**: PostgreSQL (feature flags and dependencies)
- **Audit Database**: MongoDB (operation logs)
- **ORM**: GORM for PostgreSQL operations
- **Testing**: Ginkgo + Gomega
- **Containerization**: Docker + Docker Compose

## Quick Start

### Prerequisites
- Docker and Docker Compose
- Git

### Setup

1. **Clone the repository**
   ```bash
   git clone https://github.com/ArshiAbolghasemi/dom-cobb
   cd dom-cobb
   ```

2. **Configure environment**
   ```bash
   cp .env.example .env
   # Edit .env if needed for your environment
   ```

3. **Start the services**
   ```bash
   docker-compose up --build -d
   ```

4. **Initialize the database**

   The PostgreSQL database will be automatically initialized with the required tables. If you need to manually create them:
   ```sql
   CREATE DATABASE dom_cobb;

   CREATE TABLE feature_flags (
       id SERIAL PRIMARY KEY,
       created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
       updated_at TIMESTAMP WITH TIME ZONE,
       deleted_at TIMESTAMP WITH TIME ZONE,
       "name" VARCHAR(255) NOT NULL,
       is_active BOOLEAN NOT NULL DEFAULT FALSE
   );

   CREATE UNIQUE INDEX idx_feature_flags_name ON feature_flags (name);
   CREATE INDEX idx_feature_flags_deleted_at ON feature_flags (deleted_at);

   CREATE TABLE flag_dependencies (
       flag_id BIGINT NOT NULL,
       depends_on_flag_id BIGINT NOT NULL,
       created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
       PRIMARY KEY (flag_id, depends_on_flag_id),
       CONSTRAINT fk_flag_dependencies_flag_id
           FOREIGN KEY (flag_id) REFERENCES feature_flags (id) ON DELETE CASCADE,
       CONSTRAINT fk_flag_dependencies_depends_on_flag_id
           FOREIGN KEY (depends_on_flag_id) REFERENCES feature_flags (id) ON DELETE CASCADE
   );
   ```

## Testing

Run the complete test suite:

```bash
docker build -t dom-cobb-tests -f Dockerfile.test . && \
docker run --rm \
  -v $(pwd):/app \
  -v go-cache:/go/pkg/mod \
  -v go-build-cache:/root/.cache/go-build \
  -w /app \
  dom-cobb-tests
```

## API Documentation

You can see dom-cobb's swagger in this url
```
http://localhost:8080/swagger/index.html
```
