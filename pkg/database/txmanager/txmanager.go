package txmanager

import (
	"context"
	"database/sql"
	"errors"

	"github.com/escoutdoor/study-platform/pkg/database"
	"github.com/escoutdoor/study-platform/pkg/database/pq"
	"github.com/escoutdoor/study-platform/pkg/errwrap"
)

type manager struct {
	db database.Transactor
}

func NewTransactionManager(db database.Transactor) database.TxManager {
	return &manager{
		db: db,
	}
}

func (m *manager) ReadCommited(ctx context.Context, fn database.Handler) error {
	return m.transaction(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted}, fn)
}

func (m *manager) transaction(ctx context.Context, opts *sql.TxOptions, fn database.Handler) (err error) {
	tx, ok := ctx.Value(pq.TxKey).(*sql.Tx)
	if ok && tx != nil {
		return fn(ctx)
	}

	tx, err = m.db.BeginTx(ctx, opts)
	if err != nil {
		return errwrap.Wrap("can't begin transaction", err)
	}

	ctx = pq.MakeContextTx(ctx, tx)

	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
			panic(r)
		}

		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				err = errors.Join(err, errwrap.Wrap("transaction rollback failed", rollbackErr))
			}
			return
		}

		if commitErr := tx.Commit(); commitErr != nil {
			err = errwrap.Wrap("transaction commit failed", commitErr)
		}
	}()

	err = fn(ctx)
	return err
}
