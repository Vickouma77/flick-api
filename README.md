# flick-api

`flick-api` is a Go-based API service.

## Directory Layout

```text
.
├── bin/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
├── migrations/
├── remote/
├── go.mod
└── Makefile
```

## Directory and File Descriptions

- **bin/**: Contains compiled application binaries for deployment purposes.
- **cmd/api/**: Contains application-specific entrypoint and server code (for example, starting the server and handling HTTP requests and responses).
- **internal/**: Contains packages used by the API that are not intended for external use.
- **migrations/**: Contains SQL migration files for the database.
- **remote/**: Contains configuration files and scripts for production servers.
- **Makefile**: Contains commands for automating common tasks such as code auditing, testing, and building the application.