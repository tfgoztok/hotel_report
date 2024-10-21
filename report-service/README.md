# Report Service

## Overview

The Report Service is a .NET Core-based microservice responsible for generating, storing, and retrieving statistical reports about hotels. It communicates with the Hotel Service via GraphQL and uses RabbitMQ for asynchronous report generation.

## Features

- Asynchronous report generation based on location
- Listing all generated reports
- Retrieving detailed report information

## Tech Stack

- .NET Core 7.0
- MongoDB
- RabbitMQ
- GraphQL (for communication with Hotel Service)
- Docker

## Project Structure

```
report-service/
├── src/
│   └── ReportService/
│       ├── Controllers/
│       ├── Models/
│       ├── Services/
│       ├── Interfaces/
│       ├── Data/
│       │   └── Repositories/
│       ├── Program.cs
│       └── appsettings.json
├── tests/
│   └── ReportService.Tests/
├── Dockerfile
├── docker-compose.yml
└── README.md
```

## Getting Started

### Prerequisites

- .NET Core 7.0 SDK
- Docker and Docker Compose (for running with dependencies)

### Running Locally

1. Set up the environment variables (see `appsettings.json` for required variables)

2. Run MongoDB and RabbitMQ (you can use the provided `docker-compose.yml`)

3. Navigate to the project directory:
   ```
   cd src/ReportService
   ```

4. Run the service:
   ```
   dotnet run
   ```

The service will be available at `http://localhost:5000`.

### Running with Docker

1. Build the Docker image:
   ```
   docker build -t report-service .
   ```

2. Run the service along with its dependencies:
   ```
   docker-compose up -d
   ```

## API Endpoints

- `GET /reports` - List all reports
- `GET /reports/{id}` - Get detailed report information

## Testing

To run the tests:

```
dotnet test
```

## Logging

This service uses Serilog for structured logging. Logs are sent to Elasticsearch and can be viewed using Kibana.