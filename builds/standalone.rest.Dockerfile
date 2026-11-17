# Runs the migrations, then the REST server, in one image for local development and tests.
# A production deployment runs the rest.Dockerfile image after the migrations job instead.
FROM docker.io/library/golang:1.27.4-alpine AS builder

ENV CGO_ENABLED=0

WORKDIR /app

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# The whole source trees: .dockerignore drops tests and docs, and the build cache mount
# recompiles only the packages that changed.
COPY ./cmd ./cmd
COPY ./internal ./internal

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -ldflags="-s -w" -trimpath -o /rest ./cmd/rest/ && \
    go build -ldflags="-s -w" -trimpath -o /migrations ./cmd/migrations/

FROM docker.io/library/alpine:3.24.4

COPY --from=builder /rest /rest
COPY --from=builder /migrations /migrations

# The healthcheck uses the wget of Alpine's BusyBox.
HEALTHCHECK --interval=1s --timeout=5s --retries=10 --start-period=1s \
  CMD ["wget", "-qO", "/dev/null", "http://localhost:8080/v0/ping"]

ENV REST_PORT=8080

# A local or test stack has no service-token issuer: the servers take every call.
ENV SERVICE_ISSUER_OPTIONAL=true

EXPOSE 8080

CMD ["sh", "-c", "/migrations && /rest"]
