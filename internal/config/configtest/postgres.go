package configtest

import (
	"github.com/uptrace/bun/driver/pgdriver"

	postgrespresets "github.com/a-novel-kit/golib/postgres/presets"

	"github.com/a-novel/service-template/internal/config/env"
)

// PostgresPreset is the PostgreSQL configuration used by integration tests
// (DAO tests and health/status handler tests). It points at the same database
// as the production preset; the test harness (postgres.RunTransactionalTest)
// isolates each test in a rolled-back transaction.
//
// This lives in a dedicated, test-only package — not in the production `config`
// package — so it is never compiled into a service binary. See the write-go-tests
// skill ("Cross-Package Test Fixtures").
var PostgresPreset = postgrespresets.NewDefault(pgdriver.WithDSN(env.PostgresDsn))
