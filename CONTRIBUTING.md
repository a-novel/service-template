# Contributing to service-template

For platform-wide setup (Go, Node, Podman, the `a-novel` CLI) and the day-to-day `a-novel` / `pnpm` commands, see the [developer onboarding guide](https://github.com/a-novel-kit/.github/blob/master/README.md). This file documents what is specific to the template.

For deployment, configuration, and client-package integration — and for how to fork this template into a real service — read the [README](./README.md) first.

---

## Quick local interactions

Once the service is up (`a-novel run start service-template/rest` and/or `.../grpc`), the REST server listens on `${SERVICE_TEMPLATE_REST_PORT}` and the gRPC server on `${SERVICE_TEMPLATE_GRPC_PORT}`. Both ports are allocated by the `a-novel` daemon; inject them into your shell with `eval "$(a-novel run env service-template)"`.

### Health

```bash
# REST: liveness
curl http://localhost:${SERVICE_TEMPLATE_REST_PORT}/ping

# REST: dependency check (Postgres ping)
curl http://localhost:${SERVICE_TEMPLATE_REST_PORT}/healthcheck

# gRPC: dependency check
grpcurl -plaintext localhost:${SERVICE_TEMPLATE_GRPC_PORT} StatusService/Status
```

### Item CRUD

```bash
# Create — REST then gRPC
curl -X POST "http://localhost:${SERVICE_TEMPLATE_REST_PORT}/items" \
  -H "Content-Type: application/json" \
  -d '{"name": "My Item", "description": "An optional description."}'
grpcurl -plaintext -d '{"name": "My Item", "description": "An optional description."}' \
  localhost:${SERVICE_TEMPLATE_GRPC_PORT} ItemCreateService/ItemCreate

# List (paginated)
curl "http://localhost:${SERVICE_TEMPLATE_REST_PORT}/items?limit=10&offset=0"
grpcurl -plaintext -d '{"limit": 10, "offset": 0}' localhost:${SERVICE_TEMPLATE_GRPC_PORT} ItemListService/ItemList

# Get / Update / Delete by ID
curl "http://localhost:${SERVICE_TEMPLATE_REST_PORT}/item?id=<item-uuid>"
curl -X PUT "http://localhost:${SERVICE_TEMPLATE_REST_PORT}/item" \
  -H "Content-Type: application/json" \
  -d '{"id": "<item-uuid>", "name": "Updated Name", "description": "Updated description."}'
curl -X DELETE "http://localhost:${SERVICE_TEMPLATE_REST_PORT}/item?id=<item-uuid>"
grpcurl -plaintext -d '{"id": "<item-uuid>"}' localhost:${SERVICE_TEMPLATE_GRPC_PORT} ItemGetService/ItemGet
grpcurl -plaintext -d '{"id": "<item-uuid>", "name": "Updated Name", "description": "Updated description."}' \
  localhost:${SERVICE_TEMPLATE_GRPC_PORT} ItemUpdateService/ItemUpdate
grpcurl -plaintext -d '{"id": "<item-uuid>"}' localhost:${SERVICE_TEMPLATE_GRPC_PORT} ItemDeleteService/ItemDelete
```

---

## Service-specific concepts

This is a **template**: the only domain object is the placeholder `item` resource, and the sections below document it as a worked example a fork replaces with its own resource (see [Using this template](./README.md#using-this-template)).

### The `item` resource

`item` is a named entity with an optional description, exposed through full CRUD. It is intentionally trivial — its purpose is to show one resource wired through every layer so the wiring is what you copy, not the domain logic.

### Layer architecture

The code follows the standard clean-architecture split; the import direction is one-way, handler → service → DAO.

| Layer   | Package              | Responsibility                                          | Example files                             |
| ------- | -------------------- | ------------------------------------------------------- | ----------------------------------------- |
| DAO     | `internal/dao/`      | Raw Postgres queries via the bun ORM (one `.go`/`.sql`) | `pg.item.go`, `pg.itemCreate.{go,sql}`    |
| Service | `internal/services/` | Validation, business logic, UUID generation             | `item.go`, `itemCreate.go`, `validate.go` |
| Handler | `internal/handlers/` | REST + gRPC request/response mapping                    | `http.item*.go`, `grpc.item*.go`          |

Validation uses [`go-playground/validator`](https://github.com/go-playground/validator); `internal/services/validate.go` registers one custom `notblank` tag — add your own there as needed.

### Surfaces

| Surface           | Audience                       | Operations                                                    | Spec                                                 |
| ----------------- | ------------------------------ | ------------------------------------------------------------- | ---------------------------------------------------- |
| gRPC (`cmd/grpc`) | Internal, private network only | `StatusService`, `Item{Create,Get,List,Update,Delete}Service` | [`internal/models/proto/`](./internal/models/proto/) |
| REST (`cmd/rest`) | Public, unauthenticated        | `/ping`, `/healthcheck`, `/items`, `/item`                    | [`openapi.yaml`](./openapi.yaml)                     |

The gRPC server implements no application-layer authentication — access control on that surface is enforced entirely by deployment infrastructure (network policy, ingress, service mesh). The REST reference is published at [a-novel.github.io/service-template](https://a-novel.github.io/service-template) (`openapi.yaml` + the `openapi.html` Scalar viewer).

---

## Questions?

- Open an issue at https://github.com/a-novel/service-template/issues
- Check existing issues for similar problems
- Include relevant logs and environment details
