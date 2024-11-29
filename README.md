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

-   `go run cmd/cli/main.go authorize -p <proxy_wallet_address> -b <bank_address>`
-   `go run cmd/cli/main.go balance -b <bank_address>`
-   `go run cmd/cli/main.go deposit -a <amount> -b <bank_address>`
-   `go run cmd/cli/main.go refund -b <bank_address>`
-   `go run cmd/cli/main.go withdraw -a <amount> -d <destination> -b <bank_address>`
-   `go run cmd/cli/main.go banks -p <proxy_address>`
-   `go run cmd/cli/main.go retrieval -p <proxy_address> -i <piece_cid> -a <authorization>`

## Service/Bank

To run the Bank: `go run cmd/bank/main.go`

HTTP server API featuring the following endpoints:

-   GET `/api/v1/healthcheck`: healthcheck to verify if the server is properly running
-   POST `/api/v1/register`: registers a proxy on the bank
-   POST `/api/v1/deposit`: client deposits FIL funds on the bank
-   POST `/api/v1/withdraw`: client withdraws FIL funds from the bank
-   GET `/api/v1/balance`: checks client's balance
-   POST `/api/v1/authorize`: authorizes transaction
-   GET `/api/v1/refund`: client refunds all the expired FIL funds on escrow
-   POST `/api/v1/redeem`: proxy redeems funds of transaction
-   POST `/api/v1/verify`: proxy verifies an authorization

### Migrations

Migrations are managed by [go-migrate](https://github.com/golang-migrate/migrate#cli-usage)

To run migrations:
`migrate -path=./bank/postgres/migrations -database=$DSN up`

-   `$DSN` should contain your database data source string

## Proxy

To run the Proxy: `go run cmd/proxy/main.go`

HTTP server API featuring the following endpoints:

-   GET `/api/v1/healthcheck`: healthcheck to verify if the server is properly running
-   GET `/api/v1/banks`: show the banks that the proxy is registered with
-   GET `/api/v1/fetch/{piece_cid}`: to request a file retrieval to booster-http, given a `piece-cid`
