> scaffolding based on GIN + WIRE (DI).

## Features

- Follow the `RESTful API` design specification
- Use `Wire` to resolve dependencies between modules
- Provides rich `Gin` middlewares (JWTAuth,CORS,RequestLogger,RequestRateLimiter,TraceID,CasbinEnforce,Recover,GZIP)

## Dependent Tools

```bash
go get -u github.com/cosmtrek/air
go get -u github.com/google/wire/cmd/wire
```

- [air](https://github.com/cosmtrek/air) -- Live reload for Go apps
- [wire](https://github.com/google/wire) -- Compile-time Dependency Injection for Go

## Dependent Library

- [Gin](https://gin-gonic.com/) -- The fastest full-featured web framework for Go.
- [Wire](https://github.com/google/wire) -- Compile-time Dependency Injection for Go

## Getting Started

```bash
$ go run -ldflags "-X main.VERSION=$(RELEASE_TAG)" ./main.go --config ./configs/config.local.toml --version "0.0.0"

# use Makefile:
$ make start
```

#### Auto-Rebuild

```bash
# use air for development to auto-rebuild the project when something changed
$ make dev
```

#### Escape Analysis

```bash
# default target is cmd/relajet-admin/main.go
$ make analysis

# specify target
$ make analysis TARGET=<path_to_file>
```

### Use `wire` to generate dependency injection

```bash
$ wire gen ./internal/app

# Or use Makefile:
$ make wire
```
