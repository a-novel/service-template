# Contributing to service-template

This document covers what is specific to the template. The architecture, layers and conventions every service shares are in the [service & architecture concepts](https://github.com/a-novel/.github/blob/master/CONTRIBUTING.md). Platform setup and daily commands are in the [developer guide](https://a-novel-kit.github.io/stack/).

The template is a fork target: an `item` resource implements the common service contracts end to end, with no feature of its own. The [README](./README.md#using-this-template) covers how to fork it.

---

## Running it locally

The `a-novel` daemon skips the template, so start it with Podman Compose. Set the host ports on the first lines.

```bash
POSTGRES_PORT=5432
REST_PORT=8080
GRPC_PORT=8081

export POSTGRES_PORT REST_PORT GRPC_PORT
podman compose -p service-template -f builds/podman-compose.yaml --profile rest --profile grpc up -d --build
```

A fork is a service the daemon runs, so it starts with `a-novel run start <service>/rest` instead.

Check that both servers answer:

```bash
curl "http://localhost:$REST_PORT/v0/ping"
curl "http://localhost:$REST_PORT/v0/healthcheck"
grpcurl -plaintext "localhost:$GRPC_PORT" anovel.template.v0.StatusService/Status
```

The `item` routes and RPCs are the placeholder a fork replaces. Their shapes are in [`openapi.yaml`](./openapi.yaml) and [`internal/models/proto/anovel/template/v0/`](./internal/models/proto/anovel/template/v0/).

Stop the stack and drop its database when you are done:

```bash
podman compose -p service-template -f builds/podman-compose.yaml --profile rest --profile grpc down --volumes
```

---

## Transactions

Writes that must land together run in a `transaction.Transactor`. The service that needs one takes it as a constructor dependency, and `cmd` injects it. The transactor names no database, so business code says "these writes are one unit" without knowing what stores them:

```go
// internal/core
type SomeService struct {
	dao        SomeServiceDao
	transactor transaction.Transactor
}

err := service.transactor.WithinTx(ctx, func(ctx context.Context) error {
	// every data-access call made with this ctx is part of one transaction
})

// cmd
service := core.NewSomeService(daoSomething, postgres.NewTransactor(nil))
```

Each `item` operation writes once, so none of them takes a transactor.

**Pass the callback's `ctx` down, never the outer one.** A data-access object resolves its database handle from the context, and the transaction lives on the context the callback receives. An inner call given the outer context runs on the connection pool and commits on its own, while the block still reports success.

The shared library's documentation is the contract for two more rules:

- **Never call an external service inside `WithinTx`.** An open transaction holds a pooled connection for its whole lifetime. Holding one through a third-party call exhausts the pool and blocks vacuuming. Persist what the call needs, close the transaction, make the call, then open a new transaction to record the result. `postgres.InTx(ctx)` reports whether a transaction is open, so a data-access object that makes an outbound call can refuse to run inside one.
- **A nested `WithinTx` joins the transaction in progress.** A rollback anywhere then discards the whole outermost unit of work, including work the outer caller believed was safe. A nested call also ignores its own `sql.TxOptions`, so an operation that needs a specific isolation level must be the outermost transaction.

Unit-test a service that takes a transactor with `transactiontest.NewTransactor`, which runs the callback inline. `NewFailingTransactor` covers the path where the transaction never opens: asserting that no dependency is reached proves the writes sit inside the transaction. A test that needs a real rollback needs a real database. Use `postgrestest.RunDBTest` for it, since the passthrough transaction of `postgrestest.RunTransactionalTest` cannot tell a working transactor from a broken one.

---

## Questions?

[Open an issue](https://github.com/a-novel/service-template/issues) with logs and environment details.
