# This image exposes our app as a gRPC server.
#
# It requires a patched database instance to run properly.
FROM docker.io/library/golang:1.26.4-alpine AS builder

# Static binary: no libc dependency, so it runs on a bare alpine base.
ENV CGO_ENABLED=0

WORKDIR /app

# Resolve modules off go.mod/go.sum alone, before the source, so the download
# layer survives any source-only edit. The cache mount persists the module cache
# across builds.
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Install grpcurl (healthcheck) before the source COPY so its cached compile
# survives source-only rebuilds. Pinned for reproducibility.
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    GOBIN=/usr/local/bin go install github.com/fullstorydev/grpcurl/cmd/grpcurl@v1.9.3

COPY ./cmd/grpc ./cmd/grpc
COPY ./internal/handlers ./internal/handlers
COPY ./internal/dao ./internal/dao
COPY ./internal/core ./internal/core
COPY ./internal/models ./internal/models
COPY ./internal/config ./internal/config

# Cache mounts persist the module + build cache across builds. -ldflags="-s -w"
# strips the symbol table + DWARF (smaller binary); -trimpath drops absolute paths.
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -ldflags="-s -w" -trimpath -o /grpc ./cmd/grpc/

FROM docker.io/library/alpine:3.24.1

COPY --from=builder /grpc /grpc
COPY --from=builder /usr/local/bin/grpcurl /usr/local/bin/grpcurl

HEALTHCHECK --interval=1s --timeout=5s --retries=10 --start-period=1s \
  CMD grpcurl --plaintext -d '' localhost:8080 grpc.health.v1.Health/Check || exit 1

ENV GRPC_PORT=8080

# gRPC port.
EXPOSE 8080
# TLS port.
EXPOSE 443

CMD ["/grpc"]
