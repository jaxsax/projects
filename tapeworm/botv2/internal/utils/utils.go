package utils

import (
	"context"
	"database/sql"
	"fmt"
)

type transactionContextKey struct{}
type requestIDContextKey struct{}

var ErrTransactionNotfound = fmt.Errorf("transaction not found")
var ErrRequestIDNotFound = fmt.Errorf("request id not found")

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

func SetRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDContextKey{}, id)
}

func MustGetRequestID(ctx context.Context) string {
	v, ok := ctx.Value(requestIDContextKey{}).(string)
	if !ok {
		panic(ErrRequestIDNotFound)
	}

	return v
}
