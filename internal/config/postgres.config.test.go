package config

import (
	"github.com/uptrace/bun/driver/pgdriver"

	postgrespresets "github.com/a-novel-kit/golib/postgres/presets"

	"github.com/a-novel/service-template/internal/config/env"
)

var PostgresPresetTest = postgrespresets.NewDefault(pgdriver.WithDSN(env.PostgresDsn))
