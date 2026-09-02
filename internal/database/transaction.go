package database

import (
	"context"

	db "liftwork/internal/database/sqlc"

	"github.com/jackc/pgx/v5"
)

type Transactor interface {
	db.DBTX
	Begin(ctx context.Context) (pgx.Tx, error)
}
