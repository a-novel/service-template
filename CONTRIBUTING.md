# Contributing to service-template

This document covers what is specific to the template. The [taxonomy](https://github.com/a-novel/.github/blob/master/docs/taxonomy.md#inside-a-service-layers-and-contracts) explains the layers every service shares. Platform setup and daily commands are in the [developer guide](https://a-novel-kit.github.io/stack/).

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

The golib [`transaction` package](https://pkg.go.dev/github.com/a-novel-kit/golib/transaction) holds two more rules. Call no external service inside `WithinTx`. A nested `WithinTx` joins the transaction in progress, so a rollback anywhere discards the whole unit of work. `postgres.InTx(ctx)` reports whether a transaction is open, so a data-access object that makes an outbound call can refuse to run inside one.

Unit-test a service that takes a transactor with `transactiontest.NewTransactor`, which runs the callback inline. `NewFailingTransactor` covers the path where the transaction never opens: asserting that no dependency is reached proves the writes sit inside the transaction. A test that needs a real rollback needs a real database. Use `postgrestest.RunDBTest` for it, since the passthrough transaction of `postgrestest.RunTransactionalTest` cannot tell a working transactor from a broken one.

---

## Code license

The copyrightable source code that the Licensor owns in this repository is released under [RPL-1.5](LICENSE). [NOTICE.md](NOTICE.md) names the Licensor and holds the License Notice and the licensing history.

## Contributor agreement

Before submitting a contribution for inclusion in an official A-Novel repository, review the [Contributor License Agreement](CLA.md). The Project Owner may require a signed individual or corporate agreement before accepting the contribution. A pull request, a checkbox or a signed-off commit does not replace that signature.

If you developed the contribution for an employer or a client, or it includes third-party material, make sure you may grant the rights in the agreement, and disclose any third-party license terms.

The agreement gives the Project Owner additional rights over accepted contributions, including proprietary relicensing where lawful. It changes neither the RPL-1.5 license of public releases nor anything for the maintainers of independent forks who have not signed it.

## Submitting changes

1. Open an issue or pull request describing the proposed change.
2. Identify any incorporated third-party material and its license notices.
3. When asked, sign the agreement and identify your GitHub account, the repository and the pull request or commit it covers.
4. Keep existing copyright and license notices intact.

The Project Owner records signatures and matches them to accepted pull requests or commits. GitHub neither executes nor validates the agreement.

---

## Questions?

[Open an issue](https://github.com/a-novel/service-template/issues) with logs and environment details.
