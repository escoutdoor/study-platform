package database

import (
	"context"
	"database/sql"
)

type DB interface {
	ExecContext(ctx context.Context, q Query, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, q Query, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, q Query, args ...any) *sql.Row

	Ping(ctx context.Context) error
	Conn() *sql.DB
	Close()

	Transactor
}

type Query struct {
	Name string
	Sql  string
}

type Transactor interface {
	BeginTx(ctx context.Context, txOpts *sql.TxOptions) (*sql.Tx, error)
}

type TxManager interface {
	ReadCommited(ctx context.Context, fn Handler) error
}

type Handler func(ctx context.Context) error
