package tx

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"go.uber.org/multierr"
	"platform/pkg/ctxlog"
	"platform/pkg/db/dbtx"
)

type Manager interface {
	Execute(ctx context.Context, txFunc func(context.Context) error, opts ...sql.TxOptions) error
}

type ctxMarker struct {
}

var ctxKey = ctxMarker{}

func Tx(ctx context.Context) *sqlx.Tx {
	tx, _ := ctx.Value(ctxKey).(*sqlx.Tx)
	return tx
}

func WithTx(ctx context.Context, tx *sqlx.Tx) context.Context {
	return context.WithValue(ctx, ctxKey, tx)
}

func NewManager(db dbtx.DBTX) Manager {
	return &txManager{db: db}
}

type txManager struct {
	db dbtx.DBTX
}

func (t *txManager) Execute(ctx context.Context, txFunc func(context.Context) error, opts ...sql.TxOptions) (errTx error) {
	tx, err := t.db.BeginTxx(ctx, chainOptions(opts))
	if err != nil {
		return fmt.Errorf("couldn't begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			errTx = fmt.Errorf("recovered panic in a transaction: %v", p)
			if err := tx.Rollback(); err != nil {
				ctxlog.Error(ctx, "couldn't rollback transaction in panic", ctxlog.ErrorAttr(err))
			}
		}
	}()

	err = txFunc(WithTx(ctx, tx))
	if err != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			if errors.Is(errRollback, sql.ErrTxDone) {
				ctxlog.Warn(ctx, "couldn't rollback transaction", ctxlog.ErrorAttr(errRollback))
				return fmt.Errorf("couldn't rollback transaction: %w", multierr.Append(err, errRollback))
			}
			ctxlog.Error(ctx, "couldn't rollback transaction", ctxlog.ErrorAttr(errRollback))
			return fmt.Errorf("couldn't rollback transaction: %w", multierr.Append(err, errRollback))
		}
		return err
	}
	if errCommit := tx.Commit(); errCommit != nil {
		return fmt.Errorf("couldn't commit transaction: %w", errCommit)
	}
	return nil
}

func chainOptions(opts []sql.TxOptions) *sql.TxOptions {
	var opts0 *sql.TxOptions
	if len(opts) > 0 {
		opts0 = &opts[0]

		for i := 1; i < len(opts); i++ {
			if opts[i].Isolation > opts0.Isolation {
				opts0.Isolation = opts[i].Isolation
			}
			if opts[i].ReadOnly {
				opts0.ReadOnly = true
			}
		}
	}
	return opts0
}
