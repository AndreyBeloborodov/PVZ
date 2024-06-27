package transactor

import (
	"context"
	"github.com/pkg/errors"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type QueryEngine interface {
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
}

type QueryEngineProvider interface {
	GetQueryEngine(ctx context.Context) QueryEngine // tx OR Pool
}

type TransactionManager struct {
	Pool *pgxpool.Pool
}

func NewTransactionManager(pool *pgxpool.Pool) *TransactionManager {
	return &TransactionManager{pool}
}

const key = "tx"

func (tm *TransactionManager) RunReadCommitted(ctx context.Context, fx func(ctxTX context.Context) error) error {
	tx, err := tm.Pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		return err
	}

	if err := fx(context.WithValue(ctx, key, tx)); err != nil {
		return errors.Wrap(err, tx.Rollback(ctx).Error())
	}

	if err := tx.Commit(ctx); err != nil {
		return errors.Wrap(err, tx.Rollback(ctx).Error())
	}

	return nil
}

func (tm *TransactionManager) GetQueryEngine(ctx context.Context) QueryEngine {
	tx, ok := ctx.Value(key).(QueryEngine)
	if ok && tx != nil {
		return tx
	}

	return tm.Pool
}
