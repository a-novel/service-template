// Command migrations applies the pending database schema migrations, then exits. The cmd/grpc and
// cmd/rest servers start once it completes.
package main

import (
	"context"

	"github.com/samber/lo"

	"github.com/a-novel-kit/golib/postgres"

	"github.com/a-novel/service-template/internal/config"
	"github.com/a-novel/service-template/internal/models/migrations"
)

func main() {
	ctx := context.Background()

	db := lo.Must(config.PostgresPresetMigrations.DB(ctx))
	defer config.PostgresPresetMigrations.Close()

	lo.Must0(postgres.RunMigrations(ctx, db, migrations.Migrations))
}
