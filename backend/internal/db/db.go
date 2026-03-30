package db

import (
	"database/sql"

	"github.com/rs/zerolog"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

func New(dsn string, debug bool, log zerolog.Logger) *bun.DB {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())

	db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(debug),
		bundebug.WithWriter(&logWriter{log: log}),
	))

	return db
}

type logWriter struct {
	log zerolog.Logger
}

func (w *logWriter) Write(p []byte) (n int, err error) {
	w.log.Debug().Msg(string(p))
	return len(p), nil
}
