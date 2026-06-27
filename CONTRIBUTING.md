# Contributing to service-template

Platform setup and day-to-day commands are in the [developer onboarding guide](https://github.com/a-novel-kit/.github/blob/master/README.md). This file covers what's specific to the template.

Read the [README](./README.md) first — including how to fork the template into a real service.

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

## Template-specific concepts

This is a **template**: the only domain object is the placeholder `item` resource, and the sections below document it as a worked example a fork replaces with its own resource (see [Using this template](./README.md#using-this-template)).

The vocabulary used below — **core**, **DAO**, **handler**, **server**, **job**, **API** — is platform-wide and defined once in the [service & architecture concepts](https://github.com/a-novel/.github/blob/master/CONTRIBUTING.md). This section is only the concrete shape those concepts take in the template.

### The `item` resource

`item` is a named entity with an optional description, exposed through full CRUD. It is intentionally trivial — its purpose is to show one resource wired through every layer so the wiring is what you copy, not the domain logic.

`service-template` follows the platform's standard [code structure](https://github.com/a-novel/.github/blob/master/CONTRIBUTING.md#code-structure); the `item` resource appears in each layer as:

| Layer   | `item` files                                                                           |
| ------- | -------------------------------------------------------------------------------------- |
| DAO     | `internal/dao/pg.item.go`, `pg.itemCreate.{go,sql}` (+ `Get`/`List`/`Update`/`Delete`) |
| Core    | `internal/core/item.go`, `itemCreate.go` and siblings                                  |
| Handler | `internal/handlers/http.item*.go`, `grpc.item*.go`                                     |

The core's DAO interfaces follow the `<Operation>Dao` shape — `ItemCreateDao`, `ItemListDao`, and so on. Validation lives in `internal/core/validate.go`, which registers one custom `notblank` tag — add your own there as needed.

### APIs

The long-running **server** target exposes both APIs below; the one-shot migration **job** (`jobs/migrations`) exposes none.

| API               | Audience                       | Operations                                                    | Spec                                                 |
| ----------------- | ------------------------------ | ------------------------------------------------------------- | ---------------------------------------------------- |
| gRPC (`cmd/grpc`) | Internal, private network only | `StatusService`, `Item{Create,Get,List,Update,Delete}Service` | [`internal/models/proto/`](./internal/models/proto/) |
| REST (`cmd/rest`) | Public, unauthenticated        | `/ping`, `/healthcheck`, `/items`, `/item`                    | [`openapi.yaml`](./openapi.yaml)                     |

The gRPC server implements no application-layer authentication — access control on that API is enforced entirely by deployment infrastructure (network policy, ingress, service mesh). The REST reference is published at [a-novel.github.io/service-template](https://a-novel.github.io/service-template) (`openapi.yaml` + the `openapi.html` Scalar viewer).

---

## Questions?

[Open an issue](https://github.com/a-novel/service-template/issues) — include logs and environment details.
