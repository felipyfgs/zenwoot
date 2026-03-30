package db

import (
	"context"
	"embed"

	"github.com/felipyfgs/zenwoot/backend/internal/logger"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/migrate"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

func RunMigrations(ctx context.Context, db *bun.DB) error {
	migrations := migrate.NewMigrations()
	if err := migrations.Discover(migrationsFS); err != nil {
		return err
	}
	migrator := migrate.NewMigrator(db, migrations)
	if err := migrator.Init(ctx); err != nil {
		return err
	}
	group, err := migrator.Migrate(ctx)
	if err != nil {
		return err
	}
	if group.IsZero() {
		logger.Info().Msg("migrations: already up to date")
	} else {
		logger.Info().Int("count", len(group.Migrations)).Msg("migrations applied")
	}
	return nil
}
