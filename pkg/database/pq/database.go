package pq

import (
	"context"
	"database/sql"

	"github.com/escoutdoor/study-platform/pkg/database"
	"github.com/escoutdoor/study-platform/pkg/errwrap"
	"github.com/escoutdoor/study-platform/pkg/logger"
	_ "github.com/lib/pq"
)

type key string

const TxKey key = "tx"

type db struct {
	conn *sql.DB
}

var _ database.DB = (*db)(nil)

func New(ctx context.Context, dsn string) (*db, error) {
	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, errwrap.Wrap("new database connection", err)
	}

	return &db{conn: conn}, nil
}

func (d *db) Ping(ctx context.Context) error {
	return d.conn.PingContext(ctx)
}

func (d *db) BeginTx(ctx context.Context, txOpts *sql.TxOptions) (*sql.Tx, error) {
	tx, err := d.conn.BeginTx(ctx, txOpts)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (d *db) QueryContext(ctx context.Context, q database.Query, args ...any) (*sql.Rows, error) {
	logQuery(ctx, q)

	tx, ok := ctx.Value(TxKey).(*sql.Tx)
	if ok {
		return tx.QueryContext(ctx, q.Sql, args...)
	}

	return d.conn.QueryContext(ctx, q.Sql, args...)
}

func (d *db) QueryRowContext(ctx context.Context, q database.Query, args ...any) *sql.Row {
	logQuery(ctx, q)

	tx, ok := ctx.Value(TxKey).(*sql.Tx)
	if ok {
		return tx.QueryRowContext(ctx, q.Sql, args...)
	}

	return d.conn.QueryRowContext(ctx, q.Sql, args...)
}

func (d *db) ExecContext(ctx context.Context, q database.Query, args ...any) (sql.Result, error) {
	logQuery(ctx, q)

	tx, ok := ctx.Value(TxKey).(*sql.Tx)
	if ok {
		return tx.ExecContext(ctx, q.Sql, args...)
	}

	return d.conn.ExecContext(ctx, q.Sql, args...)
}

func MakeContextTx(ctx context.Context, tx *sql.Tx) context.Context {
	return context.WithValue(ctx, TxKey, tx)
}

func (d *db) Close() {
	d.conn.Close()
}

func (d *db) Conn() *sql.DB {
	return d.conn
}

func logQuery(ctx context.Context, q database.Query) {
	logger.DebugKV(
		ctx,
		"log query",
		"sql", q.Name,
		"query", q.Sql,
	)
}
