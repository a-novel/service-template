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

## Questions?

[Open an issue](https://github.com/a-novel/service-template/issues) — include logs and environment details.
