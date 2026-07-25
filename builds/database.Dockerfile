# A postgres image with the extensions the service needs pre-loaded at build time.
#
# It does not run the service's schema migrations; run the migrations target separately.
FROM docker.io/library/postgres:18.4

ARG DEBIAN_FRONTEND=noninteractive

# Stable CI and local-development defaults. A password remains mandatory at runtime.
ENV POSTGRES_USER=postgres \
    POSTGRES_DB=postgres \
    POSTGRES_HOST_AUTH_METHOD=scram-sha-256 \
    POSTGRES_INITDB_ARGS=--auth=scram-sha-256

# ======================================================================================================================
# Prepare extension scripts.
# ======================================================================================================================
# Runs once on an empty data directory, to create the extensions.
COPY ./builds/database.sql /docker-entrypoint-initdb.d/init.sql

# ======================================================================================================================
# Finish setup.
# ======================================================================================================================
EXPOSE 5432

# Postgres does not provide a healthcheck by default.
HEALTHCHECK --interval=1s --timeout=5s --retries=10 --start-period=1s \
  CMD pg_isready || exit 1
