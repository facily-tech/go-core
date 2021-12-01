# Telemetry

Telemetry is a lib to simplify Telemetry Agents instalations.

## Installation

```sh
go get github.com/facily-tech/go-core/telemetry
```

## Implementation Examples

1. Implement it instance; (Ex: https://github.com/facily-tech/go-scaffold/blob/main/internal/container/container.go)
2. Implement it's close function into Main; (Ex: https://github.com/facily-tech/go-scaffold/blob/main/cmd/api/main.go)
3. Implement the traces. (Ex: https://github.com/facily-tech/go-scaffold/blob/main/internal/api/handler.go)