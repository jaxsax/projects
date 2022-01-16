package web

import (
	"context"
	"database/sql"
)

type transactionContextKey struct{}

func SetTransaction(ctx context.Context, tx *sql.Tx) context.Context {
	return context.WithValue(ctx, transactionContextKey{}, tx)
}

func GetTx(ctx context.Context) (*sql.Tx, bool) {
	tx, ok := ctx.Value(transactionContextKey{}).(*sql.Tx)
	if !ok {
		return nil, false
	}

	return tx, true
}

func MustGetTx(ctx context.Context) *sql.Tx {
	tx, ok := GetTx(ctx)
	if !ok {
		panic("could not find transaction")
	}

	return tx
}
