# FIDL

> Currently in development

## Table of Contents

-   [Overview](#overview)
-   [CLI (Command Line Interface)](#cli)
-   [Service/Bank](#bank)
-   [Proxy](#proxy)
-   [License](#license)

## Overview

FIDL is an on-going development project consisting in three different components:

-   [CLI (Command Line Interface)](#cli): client to deposit FIL funds and request retrievals from storage providers
-   [Service/Bank](#bank): http server to act as a service/bank
    -   [Migrations](#migrations): postgres database
-   [Proxy](#proxy): http server to be installed on storage providers, to bridge the communication between clients and the service/bank

Example of configuration files for each component on the folder `etc`.

## CLI

To run the CLI: `go run cmd/cli/main.go`

Available commands:

-   `go run cmd/cli/main.go authorize -a <amount> -b <bank_address>`
-   `go run cmd/cli/main.go balance -b <bank_address>`
-   `go run cmd/cli/main.go deposit -a <amount> -b <bank_address>`
-   `go run cmd/cli/main.go refund -b <bank_address>`
-   `go run cmd/cli/main.go withdraw -a <amount> -d <destination> -b <bank_address>`

## Service/Bank

To run the Bank: `go run cmd/bank/main.go`

HTTP server API featuring the following endpoints:

-   GET `/api/v1/healthcheck`: healthcheck to verify if the server is properly running
-   POST `/api/v1/register`: registers a proxy on the bank
-   POST `/api/v1/deposit`: client deposits FIL funds on the bank
-   GET `/api/v1/withdraw`: client withdraws FIL funds from the bank
-   POST `/api/v1/balance`: checks client's balance
-   POST `/api/v1/authorize`: authorizes transaction
-   GET `/api/v1/refund`: client refunds all the expired FIL funds on escrow
-   POST `/api/v1/redeem`: proxy redeems funds of transaction
-   POST `/api/v1/verify`: proxy verifies an authorization

### Migrations

Migrations are managed by [go-migrate](https://github.com/golang-migrate/migrate#cli-usage)

To run migrations:
`migrate -path=./postgres/migrations -database=$DSN up`

-   `$DSN` should contain your database data source string

## Proxy

To run the Proxy: `go run cmd/proxy/main.go`

HTTP server API featuring the following endpoints:

-   GET `/api/v1/healthcheck`: healthcheck to verify if the server is properly running

## License

Dual-licensed under [MIT](https://github.com/subvisual/fidl/blob/main/LICENSE-MIT) + [Apache 2.0](https://github.com/subvisual/fidl/blob/main/LICENSE-APACHE)
