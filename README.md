# SuperLeap Backend Assignment

This repository contains the backend implementation for the SuperLeap assignment.

## Prerequisites

- Go 1.25+
- Docker and Docker Compose

## Getting Started

1.  **Clone the repository**
2.  **Set up environment variables:**
    Copy `.env.example` to `.env` and adjust the values if necessary.
3.  **Start the database:**
    ```bash
    docker-compose up -d
    ```
4.  **Run the application:**
    ```bash
    go run cmd/server/main.go
    ```

## Seed Data

To populate the database with sample lead data for testing, run the following command:

```bash
docker exec -i superleap-postgres psql -U postgres -d superleap < internal/database/seed.sql
```

This will insert several leads with various statuses into the `leads` table.
