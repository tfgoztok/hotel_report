# Hotel Service

## Overview

The Hotel Service is a Go-based microservice that manages hotel information, including contact details and officials. It provides a GraphQL API for communication with other services and a REST API for external clients.

## Features

- Create and delete hotels
- Add and remove hotel contact information
- List hotel officials
- Retrieve detailed hotel information
- Initiate report generation requests

## Tech Stack

- Go 1.21
- PostgreSQL
- GraphQL
- RabbitMQ
- Docker

## Project Structure

```
hotel-service/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── api/
│   │   ├── graphql/
│   │   │   └── schema.go
│   │   ├── handlers/
│   │   │   ├── contact_handler.go
│   │   │   ├── graphql_handler.go
│   │   │   ├── hotel_handler.go
│   │   │   └── report_handler.go
│   │   ├── middleware/
│   │   │   └── logging.go
│   │   └── router.go
│   ├── config/
│   │   └── config.go
│   ├── db/
│   │   ├── migrations/
│   │   │   ├── 000001_create_hotels_table.up.sql
│   │   │   └── 000002_create_contacts_table.up.sql
│   │   ├── db.go
│   │   └── migrations.go
│   ├── messaging/
│   │   └── rabbitmq.go
│   ├── models/
│   │   ├── contact.go
│   │   └── hotel.go
│   ├── repository/
│   │   ├── contact_repository.go
│   │   └── hotel_repository.go
│   └── service/
│       ├── contact_service.go
│       └── hotel_service.go
├── pkg/
│   └── logger/
│       └── logger.go
├── tests/
│   ├── integration/
│   │   └── rabbitmq_test.go
│   └── unit/
│       ├── handlers_test.go
│       ├── rabbitmq_test.go
│       ├── repositories_test.go
│       └── services_test.go
├── Dockerfile
├── docker-compose.yml
├── filebeat.yml
└── README.md
```

## Getting Started

### Prerequisites

- Go 1.21 or later
- Docker and Docker Compose (for running with dependencies)

### Running Locally

1. Set up the environment variables (see `config/config.go` for required variables)

2. Run the PostgreSQL database and RabbitMQ (you can use the provided `docker-compose.yml`)

3. Run the migrations:
   ```
   go run cmd/migrate/main.go
   ```

4. Start the service:
   ```
   go run cmd/api/main.go
   ```

The service will be available at `http://localhost:8080`.

### Running with Docker

1. Build the Docker image:
   ```
   docker build -t hotel-service .
   ```

2. Run the service along with its dependencies:
   ```
   docker-compose up -d
   ```

## API Endpoints

### REST API

- `POST /hotels` - Create a new hotel
- `DELETE /hotels/{id}` - Remove a hotel
- `POST /hotels/{id}/contacts` - Add contact information to a hotel
- `DELETE /hotels/{id}/contacts/{contactId}` - Remove contact information from a hotel
- `GET /hotels/{id}/officials` - List hotel officials
- `GET /hotels/{id}` - Get detailed hotel information
- `POST /reports/request` - Request a new report

### GraphQL API

The GraphQL endpoint is available at `/graphql`. It provides the following queries:

- `hotelsByLocation(location: String!)`: Retrieves hotels based on location
- `contactsByLocation(location: String!)`: Retrieves contacts based on location

## Testing

To run the tests:

```
go test ./...
```

For integration tests:

```
go test ./tests/integration
```

## Logging

This service uses structured logging. Logs are sent to stdout and can be collected by the ELK stack for centralized logging.