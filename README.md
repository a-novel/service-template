# Service Template

A scaffold for new A-Novel backend services: fork it, rename the Go module, and replace the example `item` resource with your own. It implements the platform's common service contracts through a dummy `item` implementation, so a fork starts from a complete, working example of every layer.

[![X (formerly Twitter) Follow](https://img.shields.io/twitter/follow/agorastoryverse)](https://twitter.com/agorastoryverse)
[![Discord](https://img.shields.io/discord/1315240114691248138?logo=discord)](https://discord.gg/rp4Qr8cA)

<hr />

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/a-novel/service-template)
![GitHub repo file or directory count](https://img.shields.io/github/directory-file-count/a-novel/service-template)
![GitHub code size in bytes](https://img.shields.io/github/languages/code-size/a-novel/service-template)

![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/a-novel/service-template/main.yaml)
[![codecov](https://codecov.io/gh/a-novel/service-template/graph/badge.svg)](https://codecov.io/gh/a-novel/service-template)

![Coverage graph](https://codecov.io/gh/a-novel/service-template/graphs/sunburst.svg)

## Using this template

This repository is a starting point for a new A-Novel backend service, not a service you deploy as is. The rest of this README demonstrates the standard service-README structure your fork inherits — using the placeholder `item` resource. Fork it, then:

**1. Rename the module and service identity.** Run `go mod edit -module github.com/a-novel/<your-service>`, then a project-wide find/replace of `service-template` → `<your-service>` (the import paths are `github.com/a-novel/service-template/internal/...`, so the same replace catches them). This also covers the badges, image names, and doc links in this README, the CI `image_name:` fields, `package.json`, and `pkg/js/rest/package.json`. Then update:

| What             | Where                                           | From → to                                          |
| ---------------- | ----------------------------------------------- | -------------------------------------------------- |
| Root package     | `generate.go` (`package …`)                     | placeholder package name → `package <yourservice>` |
| Env-var prefix   | `internal/config/env/env.go`                    | `SERVICE_TEMPLATE_ENV_PREFIX` → `<YOUR_SERVICE>_…` |
| App name default | `internal/config/env/env.go` (`AppNameDefault`) | `"service-template"` → `"<your-service>"`          |

**2. Replace the `item` resource.** Swap `item`/`Item` for your own resource across each layer, renaming and adjusting fields:

- **Migration** — `internal/models/migrations/20250306000000_items_table.{up,down}.sql` (use `date '+%Y%m%d%H%M%S'` for a fresh timestamp).
- **DAO** — `internal/dao/pg.item.go` (bun model), `pg.itemCreate.{go,sql}` and the `Get`/`List`/`Update`/`Delete` siblings, plus their `*_test.go`.
- **Core** — `internal/core/item*.go` and tests; `internal/core/validate.go` only if your fields need validators beyond the registered `notblank`.
- **Handlers** — REST `internal/handlers/http.item*.go`, gRPC `internal/handlers/grpc.item*.go`, plus tests.
- **Proto** — `internal/models/proto/item*.proto` (then `pnpm generate` to refresh `internal/handlers/protogen/`).
- **API surface** — `openapi.yaml` (+ regenerate `openapi.html`), `pkg/go/client.go`, `pkg/js/rest/src/item.ts` + `index.ts`, and `pkg/js/test/`.
- **Wiring** — constructor calls and route registrations in `cmd/rest/main.go` and `cmd/grpc/main.go`.

**3. Keep the scaffolding.** Leave these in place (adjust only as your service needs): `internal/config/`, the `cmd/*/main.go` startup shape, `internal/handlers/http.ping.go` / `http.health.go` / `http.decoder.go` / `grpc.status.go`, `internal/models/migrations/migrations.go`, `builds/`, `.github/workflows/`, the `*.mod` tool files, `renovate.json`, `buf.*`, and the prettier/pnpm config. Update only the service-naming bits of `LICENSE`, `SECURITY.md`, `CODE_OF_CONDUCT.md`, and `CONTRIBUTING.md`.

**4. Verify.** `pnpm generate` → `pnpm format` → `pnpm lint` → `a-novel test -y` (`format` and `lint` already cover Go, proto, and JS). Refresh `openapi.html`, then delete this section and finish updating this README for the new service.

The layering, naming, error, telemetry, and test conventions your fork must follow are documented in the [service & architecture concepts](https://github.com/a-novel/.github/blob/master/CONTRIBUTING.md).

## What it does

This is an example service whose only domain object is `item` — a named entity with an optional description — exposed through full CRUD. It exists to be replaced: it demonstrates the [layered architecture](https://github.com/a-novel/.github/blob/master/CONTRIBUTING.md) (DAO → core → handler), the dual REST/gRPC API, and the client packages a real service inherits.

The service ships two APIs:

- A **private gRPC API** (`cmd/grpc`) — `StatusService` plus the `Item{Create,Get,List,Update,Delete}Service` RPCs — for internal, private-network service-to-service traffic. The server implements no application-layer authentication: access control is enforced externally (network policy, ingress, service mesh).
- A **public REST API** (`cmd/rest`) — `/ping`, `/healthcheck`, and the `/items` + `/item` CRUD routes — for any HTTP client.

## Deploying

The service runs as published OCI images plus a PostgreSQL database. Both servers are stateless, so each scales to as many replicas as you need behind a load balancer; all state lives in Postgres.

> **OpenTofu modules are the planned canonical deployment path.** Until they land, deploy the images with any container orchestrator — the composition below is the reference for which images to run, how they wire together, and the environment they expect.

| Image                              | Role                                                                        |
| ---------------------------------- | --------------------------------------------------------------------------- |
| `service-template/grpc`            | Private item CRUD + status API. Internal network only.                      |
| `service-template/rest`            | Public item CRUD + health API.                                              |
| `service-template/jobs/migrations` | One-shot schema migration job; runs to completion before the servers start. |
| `service-template/database`        | Pre-tuned PostgreSQL image — or bring your own Postgres.                    |

Pin every image to the same release tag. A production deployment runs `database`, then `migrations` to completion, then any number of `grpc` and/or `rest` replicas. Provide `POSTGRES_USER`, `POSTGRES_PASSWORD`, and `POSTGRES_DATABASE` through the orchestrator; store the password as a secret.

<!-- TODO(project-docs): replace v0.0.0 with the new service's release tag -->

```yaml
services:
  postgres-template:
    image: ghcr.io/a-novel/service-template/database:v0.0.0
    networks: [api]
    environment:
      POSTGRES_PASSWORD: "${POSTGRES_PASSWORD}"
      POSTGRES_USER: "${POSTGRES_USER}"
      POSTGRES_DB: "${POSTGRES_DATABASE}"
      POSTGRES_HOST_AUTH_METHOD: scram-sha-256
      POSTGRES_INITDB_ARGS: --auth=scram-sha-256
    volumes:
      - template-postgres-data:/var/lib/postgresql/
    # Podman ignores the image's own check, so the dependents' wait needs this one.
    healthcheck:
      test: ["CMD", "pg_isready"]
      interval: 1s
      timeout: 5s
      retries: 10

  migrations-template:
    image: ghcr.io/a-novel/service-template/jobs/migrations:v0.0.0
    depends_on:
      postgres-template: { condition: service_healthy }
    environment:
      POSTGRES_HOST: postgres-template
      POSTGRES_PORT: "5432"
      POSTGRES_USER: "${POSTGRES_USER}"
      POSTGRES_PASSWORD: "${POSTGRES_PASSWORD}"
      POSTGRES_DATABASE: "${POSTGRES_DATABASE}"
      POSTGRES_TLS_ENABLED: "false"
    networks: [api]

  service-template:
    image: ghcr.io/a-novel/service-template/grpc:v0.0.0 # or .../rest:v0.0.0 for the public REST API
    ports: ["${SERVICE_TEMPLATE_GRPC_PORT}:8080"] # the container always listens on 8080; map ${SERVICE_TEMPLATE_REST_PORT} for the rest image
    depends_on:
      postgres-template: { condition: service_healthy }
      migrations-template: { condition: service_completed_successfully }
    environment:
      POSTGRES_HOST: postgres-template
      POSTGRES_PORT: "5432"
      POSTGRES_USER: "${POSTGRES_USER}"
      POSTGRES_PASSWORD: "${POSTGRES_PASSWORD}"
      POSTGRES_DATABASE: "${POSTGRES_DATABASE}"
      POSTGRES_TLS_ENABLED: "false"
    networks: [api]

networks:
  api:

volumes:
  template-postgres-data:
```

Run both servers by adding a second service that reuses the same database and migrations with the `rest` image.

### Database image

The database image packages PostgreSQL and pgBackRest on Wolfi. Its
[package manifest](./builds/database.apko.yaml) is assembled with apko during the container
build, from [locks](./builds/database.apko.lock.json) that pin every package, so an unchanged
image rebuilds to the same digest; Docker and Podman users need no additional host tools.
PostgreSQL keeps its standard `POSTGRES_*` initialization variables, `uuid-ossp` extension and
volume at `/var/lib/postgresql`. Backup scheduling and repository credentials remain the
operator's responsibility.

Use a fresh volume when moving from the Debian-based image. To keep existing data, take a
logical dump and restore it into the new database, then validate it before switching traffic.
Do not reuse a Debian data directory or assume its physical backups are portable: the libc,
collation and extension environment changed. Retain the previous image with its backups.

### Configuration

Every variable is read from the process environment. Env-var names can be globally prefixed via `SERVICE_TEMPLATE_ENV_PREFIX` (rename this when forking — see [Using this template](#using-this-template)). Set `POSTGRES_HOST` to use the discrete connection fields; `POSTGRES_DSN` remains a compatibility fallback.

| Name                   | Description                                                                                                                   | Images |
| ---------------------- | ----------------------------------------------------------------------------------------------------------------------------- | ------ |
| `POSTGRES_HOST`        | PostgreSQL hostname or IP address. Selects the discrete connection fields. **Required for new deployments.**                  | all    |
| `POSTGRES_PORT`        | PostgreSQL port. Defaults to `5432`.                                                                                          | all    |
| `POSTGRES_USER`        | PostgreSQL login role. **Required when `POSTGRES_HOST` is set.**                                                              | all    |
| `POSTGRES_PASSWORD`    | PostgreSQL login password. **Required when `POSTGRES_HOST` is set; inject it as a secret.**                                   | all    |
| `POSTGRES_DATABASE`    | PostgreSQL database name. **Required when `POSTGRES_HOST` is set.**                                                           | all    |
| `POSTGRES_TLS_ENABLED` | Encrypt the PostgreSQL connection. Defaults to `true`; disable only when another trusted boundary protects the database link. | all    |
| `POSTGRES_DSN`         | Connection-URL fallback, read only when `POSTGRES_HOST` is empty.                                                             | all    |

<details>
<summary>Optional configuration (server tuning, OpenTelemetry)</summary>

gRPC tuning (images `grpc`, `standalone-grpc`):

| Name                    | Description               | Default |
| ----------------------- | ------------------------- | ------- |
| `GRPC_TIMEOUT_SHUTDOWN` | Graceful shutdown budget. | `30s`   |

REST tuning (images `rest`, `standalone-rest`):

| Name                          | Description                                                                                             | Default                          |
| ----------------------------- | ------------------------------------------------------------------------------------------------------- | -------------------------------- |
| `REST_MAX_REQUEST_SIZE`       | Maximum request body size, in bytes.                                                                    | `16384` (16KiB)                  |
| `REST_CLIENT_IP_HEADER`       | Header a trusted proxy sets to the visitor IP. Set it only when every request comes through that proxy. | TCP peer                         |
| `REST_LIVENESS_PATH`          | Ping path, which still answers during a planned downtime and skips the rate limit.                      | `/ping`                          |
| `REST_RATE_LIMITS_ENABLED`    | Rate-limit each visitor on every route but `/ping`. Turn off only in test stacks.                       | `true`                           |
| `REST_RATE_LIMIT`             | Requests each visitor sends to each route, as `<requests>/<window>`.                                    | `5/1s`                           |
| `REST_TIMEOUT_READ`           | Read timeout.                                                                                           | `15s`                            |
| `REST_TIMEOUT_READ_HEADER`    | Header read timeout.                                                                                    | `3s`                             |
| `REST_TIMEOUT_WRITE`          | Write timeout.                                                                                          | `30s`                            |
| `REST_TIMEOUT_IDLE`           | Idle keep-alive timeout.                                                                                | `60s`                            |
| `REST_TIMEOUT_REQUEST`        | Per-request timeout.                                                                                    | `60s`                            |
| `REST_TIMEOUT_SHUTDOWN`       | Graceful shutdown budget.                                                                               | `30s`                            |
| `REST_CORS_ALLOWED_ORIGINS`   | CORS allowed origins.                                                                                   | `*`                              |
| `REST_CORS_ALLOWED_HEADERS`   | CORS allowed headers.                                                                                   | `*`                              |
| `REST_CORS_ALLOWED_METHODS`   | CORS allowed methods, comma-separated.                                                                  | `HEAD,GET,POST,PUT,PATCH,DELETE` |
| `REST_CORS_ALLOW_CREDENTIALS` | CORS allow-credentials flag.                                                                            | `false`                          |
| `REST_CORS_MAX_AGE`           | CORS max-age, in seconds.                                                                               | `3600`                           |

Database connection pool (server images). The limits are **per process**. The database's `max_connections` has to cover every replica plus the migration job; the stock `postgres` default is 100.

| Name                      | Description                               | Default |
| ------------------------- | ----------------------------------------- | ------- |
| `POSTGRES_MAX_OPEN_CONNS` | Maximum open connections to the database. | `20`    |
| `POSTGRES_MAX_IDLE_CONNS` | Maximum connections kept open while idle. | `20`    |

Planned downtime (server images):

| Name             | Description                                                                                                                                                                                                        | Default |
| ---------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | ------- |
| `DOWNTIME_START` | When a planned downtime starts, in RFC 3339. From then until it is removed, even past the announced end, the servers answer `503` / `UNAVAILABLE` except for ping and gRPC health, and start without the database. |         |

Logs and tracing — OpenTelemetry supports a stdout and a Google Cloud exporter (all server images):

| Name                | Description                                                           | Default            |
| ------------------- | --------------------------------------------------------------------- | ------------------ |
| `OTEL`              | Enable OTel tracing; the variables below pick the exporter.           | `false`            |
| `GCLOUD_PROJECT_ID` | Google Cloud project ID. When set, switches the OTel exporter to GCP. |                    |
| `APP_NAME`          | Application name attached to traces and logs.                         | `service-template` |

</details>

## Using the client packages

Two clients ship with the service. Each snippet is the **minimum viable call**; the full surface is what your editor's intellisense, [pkg.go.dev](https://pkg.go.dev/github.com/a-novel/service-template), and the [API reference](https://a-novel.github.io/service-template) are for.

- **Go** talks gRPC — use it from a backend service.
- **JavaScript / TypeScript** talks REST — use it from a frontend or Node service.

### Go (gRPC)

```bash
go get github.com/a-novel/service-template
```

```go
package main

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	servicetemplate "github.com/a-novel/service-template/pkg/go"
)

func main() {
	ctx := context.Background()

	// In production, swap insecure.NewCredentials() for a TLS or mTLS credential — the
	// server has no application-layer auth, so transport security is the only thing
	// protecting the private gRPC server from a network adversary.
	client, err := servicetemplate.NewClient(
		"service-template:8080",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	created, err := client.ItemCreate(ctx, &servicetemplate.ItemCreateRequest{
		Name:        "My Item",
		Description: "An optional description.",
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("created item %s", created.GetItem().GetId())
}
```

### JavaScript / TypeScript (REST)

The package is published to GitHub Packages, which requires a Personal Access Token with the `read:packages` scope even for public packages ([why](https://github.com/orgs/community/discussions/23386#discussioncomment-3240193)). Add to `.npmrc` (project root or `$HOME`):

```ini
@a-novel:registry=https://npm.pkg.github.com
@a-novel-kit:registry=https://npm.pkg.github.com
//npm.pkg.github.com/:_authToken=${YOUR_PERSONAL_ACCESS_TOKEN}
```

```bash
pnpm add @a-novel/service-template-rest
```

```typescript
import { TemplateApi, itemCreate, itemList } from "@a-novel/service-template-rest";

const api = new TemplateApi("http://service-template:8080");

const created = await itemCreate(api, "My Item", "An optional description.");
const items = await itemList(api, 10, 0);
```

API reference: [a-novel.github.io/service-template](https://a-novel.github.io/service-template).

## Running locally

For a throwaway instance without the dev toolchain, the **standalone** images bundle the server and migrations in one container. They run migrations on every boot — handy for a quick spin-up, unsafe under multi-replica production restarts.

```yaml
services:
  postgres-template:
    image: ghcr.io/a-novel/service-template/database:v0.0.0
    networks: [api]
    environment:
      POSTGRES_PASSWORD: postgres
      POSTGRES_USER: postgres
      POSTGRES_DB: postgres
      POSTGRES_HOST_AUTH_METHOD: scram-sha-256
      POSTGRES_INITDB_ARGS: --auth=scram-sha-256
    # Podman ignores the image's own check, so the dependents' wait needs this one.
    healthcheck:
      test: ["CMD", "pg_isready"]
      interval: 1s
      timeout: 5s
      retries: 10

  service-template:
    image: ghcr.io/a-novel/service-template/standalone-grpc:v0.0.0 # or standalone-rest
    ports: ["${SERVICE_TEMPLATE_GRPC_PORT}:8080"] # map ${SERVICE_TEMPLATE_REST_PORT} for the standalone-rest image
    depends_on:
      postgres-template: { condition: service_healthy }
    environment:
      POSTGRES_DSN: "postgres://postgres:postgres@postgres-template:5432/postgres?sslmode=disable"
    networks: [api]

networks:
  api:
```

Working on the service itself? Use the `a-novel` CLI (`a-novel run start service-template/rest`) instead — see [CONTRIBUTING](./CONTRIBUTING.md).

## Contributing

Platform setup and the day-to-day commands live in the [developer guide](https://a-novel-kit.github.io/stack/). Template-specific concepts and local interactions are in [CONTRIBUTING.md](./CONTRIBUTING.md).
