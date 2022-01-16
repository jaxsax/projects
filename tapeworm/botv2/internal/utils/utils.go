package utils

import (
	"context"
	"database/sql"
	"fmt"
)

type transactionContextKey struct{}

var ErrTransactionNotfound = fmt.Errorf("transaction not found")

func SetTransaction(ctx context.Context, tx *sql.Tx) context.Context {
	return context.WithValue(ctx, transactionContextKey{}, tx)
}

func GetTx(ctx context.Context) (*sql.Tx, error) {
	tx, ok := ctx.Value(transactionContextKey{}).(*sql.Tx)
	if !ok {
		return nil, ErrTransactionNotfound
	}

	return tx, nil
}

func MustGetTx(ctx context.Context) *sql.Tx {
	tx, err := GetTx(ctx)
	if err != nil {
		panic(err)
	}

	return tx
}
