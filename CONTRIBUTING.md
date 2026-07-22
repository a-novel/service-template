# Contributing to service-template

This file covers only what is specific to the **template**. For service-level contribution shared across every service — the architecture, the layers, the conventions — start with the [service & architecture concepts](https://github.com/a-novel/.github/blob/master/CONTRIBUTING.md). Platform setup and day-to-day commands are in the [developer onboarding guide](https://github.com/a-novel-kit/.github/blob/master/README.md).

`service-template` is a fork target: a dummy `item` resource implements the common service contracts end to end, with no real feature of its own. How to fork it — and where every `item` file lives — is in the [README](./README.md#using-this-template).

---

## Running it locally

Start a server and load its ports into your shell:

```bash
a-novel run start service-template/rest   # and/or service-template/grpc
eval "$(a-novel run env service-template)"
```

Check it is alive:

```bash
curl http://localhost:${SERVICE_TEMPLATE_REST_PORT}/ping          # REST liveness
curl http://localhost:${SERVICE_TEMPLATE_REST_PORT}/healthcheck   # REST: Postgres dependency
grpcurl -plaintext localhost:${SERVICE_TEMPLATE_GRPC_PORT} StatusService/Status   # gRPC dependency
```

The `item` CRUD routes (`/items`, `/item`, the `Item*Service` RPCs) are placeholder wiring to fork, not a feature; their request/response shapes live in [`openapi.yaml`](./openapi.yaml) and [`internal/models/proto/`](./internal/models/proto/).

---

## Transactions

Two or more writes that must land together are wrapped in a `transaction.Transactor`, taken as a constructor dependency by the service that needs one and injected in `cmd`. It names no database, so business code says "these writes are one unit" without knowing what stores them:

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

Nothing in this service uses one yet, because no operation writes twice — the `item` resource is single-write throughout. The convention is here rather than demonstrated because wrapping a single write in a transaction is noise, and a scaffold that shows it teaches every service copied from it to do the same.

**Pass the callback's `ctx` down, not the outer one.** Data-access objects resolve their database handle from the context, and the transaction is installed on the context the callback receives. An inner call given the outer context runs on the connection pool and commits on its own, while the surrounding block still reports success. That is not hypothetical: it is what a sibling service did in four operations for months, with a green build the whole time.

Two rules follow, and the shared library's documentation is the contract for both:

- **Never call an external service inside `WithinTx`.** An open transaction holds a pooled connection for its whole lifetime; pinning one for the length of a third-party call exhausts the pool and blocks vacuuming. Persist what the call needs, close the transaction, make the call, then open a new transaction to record the result. `postgres.InTx(ctx)` reports whether a transaction is open, so a data-access object that makes an outbound call can refuse rather than rely on the convention holding.
- **A nested `WithinTx` joins the transaction in progress**, so a rollback anywhere discards the whole outermost unit of work — including work the outer caller believed was already safe. Nesting is legal; it should be deliberate. A nested call also never sees its own `sql.TxOptions`, so an operation needing a specific isolation level has to be the outermost transaction.

Unit-test a service that takes a transactor with `transactiontest.NewTransactor`, which runs the callback inline, or `NewFailingTransactor` to cover the path where the unit of work never opens — asserting the dependencies are never reached is how a test proves the writes are inside the scope rather than merely near it. A test that needs a real rollback needs a real database: use `postgres.RunDBTest`, never `RunTransactionalTest`, whose passthrough transaction cannot tell a working transactor from a broken one.

---

## Questions?

[Open an issue](https://github.com/a-novel/service-template/issues) — include logs and environment details.
