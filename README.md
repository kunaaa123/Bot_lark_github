# Git-Lark Notifier

A Hexagonal Architecture based application that sends GitHub deployment notifications to Lark via webhook.

## Project Structure
```
BOT_LARK_GITHUB/
├── cmd/
│   └── main.go               # Application entry point
├── internal/
│   ├── core/                 # Business logic (Domain layer)
│   │   ├── domain/
│   │   │   └── deploy.go     # Domain models
│   │   ├── ports/
│   │   │   ├── git_port.go   # Git interface
│   │   │   └── notification_port.go # Notification interface
│   │   └── service/
│   │       └── deploy_service.go    # Core business logic
│   ├── adapters/             # Implementation layer
│   │   ├── primary/
│   │   │   └── http_handler.go      # HTTP handlers
│   │   └── secondary/
│   │       ├── github_adapter.go    # GitHub implementation
│   │       └── lark_adapter.go      # Lark implementation
│   └── config/
│       └── config.go         # Application configuration
└── go.mod                    # Go module file
```

## Architecture Overview

This project follows the Hexagonal Architecture (Ports and Adapters) pattern:

- **Core**: Contains the business logic and domain models
- **Ports**: Defines interfaces for external interactions
- **Adapters**: Implements the interfaces defined in ports

### Components

1. **Domain Layer** (`internal/core/domain`)
   - Contains business entities and logic

2. **Ports** (`internal/core/ports`)
   - Defines interfaces for external services
   - `GitPort`: Interface for Git operations
   - `NotificationPort`: Interface for notification services

3. **Adapters** (`internal/adapters`)
   - Primary: HTTP handlers for incoming requests
   - Secondary: Implementations for GitHub and Lark services

## Setup and Configuration

1. Clone the repository
2. Configure the application in `config/config.go`
3. Build and run:
```bash
go build ./cmd/main.go
./main
```

## Requirements

- Go 1.16 or later
- GitHub webhook configuration
- Lark bot webhook URL

## License

MIT License