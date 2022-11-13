package repository

import (
	"context"

	"github.com/volatiletech/sqlboiler/v4/boil"
)

// Transaction
type Transaction interface {
	Commit() error
	Rollback() error
	Executor() boil.ContextExecutor
}

// TxCreator contains all functions required to manage Room objects and their state
type TxCreator interface {
	CreateTransaction(ctx context.Context) (Transaction, error)
}
