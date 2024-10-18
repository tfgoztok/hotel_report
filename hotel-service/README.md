hotel-service/
│
├── cmd/
│   └── api/
│       └── main.go
│
├── internal/
│   ├── api/
│   │   ├── handlers/
│   │   │   ├── hotel_handler.go
│   │   │   └── contact_handler.go
│   │   ├── middleware/
│   │   │   └── logging.go
│   │   └── router.go
│   │
│   ├── config/
│   │   └── config.go
│   │
│   ├── db/
│   │   ├── migrations/
│   │   │   ├── 000001_create_hotels_table.up.sql
│   │   │   ├── 000002_create_contacts_table.up.sql
│   │   └── db.go
│   │
│   ├── models/
│   │   ├── hotel.go
│   │   └── contact.go
│   │
│   ├── repository/
│   │   ├── hotel_repository.go
│   │   └── contact_repository.go
│   │
│   └── service/
│       ├── hotel_service.go
│       └── contact_service.go
│
├── pkg/
│   └── logger/
│       └── logger.go
│
├── tests/
│   ├── integration/
│   │   └── api_test.go
│   └── unit/
│       ├── handlers_test.go
│       ├── services_test.go
│       └── repositories_test.go
│
├── Dockerfile
├── Dockerfile.test
├── docker-compose.yml
├── go.mod
├── go.sum
└── README.md