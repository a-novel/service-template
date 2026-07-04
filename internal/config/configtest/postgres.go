// Package configtest holds configuration presets shared across the service's
// integration tests.
package configtest

import (
	"github.com/uptrace/bun/driver/pgdriver"

	postgrespresets "github.com/a-novel-kit/golib/postgres/presets"

	"github.com/a-novel/service-template/internal/config/env"
)

// PostgresPreset is the PostgreSQL configuration shared by integration tests,
// such as the DAO tests. It targets the same database as the production preset;
// the transactional harness rolls back each test, so tests never observe each
// other's writes.
//
// It lives in a regular (non-_test.go) file because Go excludes _test.go files
// from a package's exported surface, and other packages' tests need to import
// this preset. Production code never imports configtest; that boundary is a
// convention enforced in review rather than by the compiler.
var PostgresPreset = postgrespresets.NewDefault(pgdriver.WithDSN(env.PostgresDsn))
