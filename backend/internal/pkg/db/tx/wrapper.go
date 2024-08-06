package tx

import (
	"context"
	"database/sql"

	"platform/internal/pkg/db/dbtx"

	"github.com/jmoiron/sqlx"
)

func Wrap(dbtx dbtx.DBTX) dbtx.DBTX {
	return &txWrapper{
		db: dbtx,
	}
}

type txWrapper struct {
	db dbtx.DBTX
}

func (t txWrapper) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	tx := Tx(ctx)
	if tx == nil {
		return t.db.ExecContext(ctx, query, args...)
	}
	return tx.ExecContext(ctx, query, args...)
}

func (t txWrapper) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	tx := Tx(ctx)
	if tx == nil {
		return t.db.QueryContext(ctx, query, args...)
	}
	return tx.QueryContext(ctx, query, args...)
}

func (t txWrapper) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	tx := Tx(ctx)
	if tx == nil {
		return t.db.QueryRowContext(ctx, query, args...)
	}
	return tx.QueryRowContext(ctx, query, args...)
}

func (t txWrapper) PrepareContext(ctx context.Context, s string) (*sql.Stmt, error) {
	tx := Tx(ctx)
	if tx == nil {
		return t.db.PrepareContext(ctx, s)
	}
	return tx.PrepareContext(ctx, s)
}

func (t txWrapper) BeginTxx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error) {
	return t.db.BeginTxx(ctx, opts)
}

func (t txWrapper) QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error) {
	tx := Tx(ctx)
	if tx == nil {
		return t.db.QueryxContext(ctx, query, args...)
	}
	return tx.QueryxContext(ctx, query, args...)
}

func (t txWrapper) QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row {
	tx := Tx(ctx)
	if tx == nil {
		return t.db.QueryRowxContext(ctx, query, args...)
	}
	return tx.QueryRowxContext(ctx, query, args...)
}

func (t txWrapper) Rebind(query string) string {
	return t.db.Rebind(query)
}

func (t txWrapper) NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error) {
	tx := Tx(ctx)
	if tx == nil {
		return t.db.NamedExecContext(ctx, query, arg)
	}
	return tx.NamedExecContext(ctx, query, arg)
}

func (t txWrapper) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	tx := Tx(ctx)
	if tx == nil {
		return t.db.SelectContext(ctx, dest, query, args...)
	}
	return tx.SelectContext(ctx, dest, query, args...)
}

func (t txWrapper) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	tx := Tx(ctx)
	if tx == nil {
		return t.db.GetContext(ctx, dest, query, args...)
	}
	return tx.GetContext(ctx, dest, query, args...)
}

func (t txWrapper) PrepareNamedContext(ctx context.Context, query string) (*sqlx.NamedStmt, error) {
	tx := Tx(ctx)
	if tx == nil {
		return t.db.PrepareNamedContext(ctx, query)
	}
	return tx.PrepareNamedContext(ctx, query)
}

func (t txWrapper) NamedQueryContext(ctx context.Context, query string, arg interface{}) (*sqlx.Rows, error) {
	tx := Tx(ctx)
	if tx == nil {
		return t.db.NamedQueryContext(ctx, query, arg)
	}
	return sqlx.NamedQueryContext(ctx, tx, query, arg)
}
