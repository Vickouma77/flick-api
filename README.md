# flick-api

.
├── bin
├── cmd
│
└── api
│   └── main.go
|
├── internal
|
├── migrations
|
├── remote
|
├── go.mod
└── Makefile

**bin**: Contains compiled application binaries for deployment purposes.
**cmd/api**: Contains application specific code, i.e running server, reading and writing HTTP requests etc
**internal**: Contains packages used by the API that are not intended to be used by external applications.
**migrations**: SQL migration files for the database.
**remote**: Contains configuration files and scripts for production servers
**Makefile**: Contains files for automating common administrative tasks i.e auditing Go code, running tests, building the application, etc.