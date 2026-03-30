package repo

import (
	"context"

	"github.com/uptrace/bun"
)

type BaseRepo[T any] struct {
	db *bun.DB
}

func NewBaseRepo[T any](db *bun.DB) *BaseRepo[T] {
	return &BaseRepo[T]{db: db}
}

func (r *BaseRepo[T]) WithTenant(ctx context.Context, accountID int64) *bun.SelectQuery {
	var m T
	return r.db.NewSelect().Model(&m).Where(`"account_id" = ?`, accountID)
}

func (r *BaseRepo[T]) DB() *bun.DB { return r.db }
