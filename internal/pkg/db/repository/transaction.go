package repository

import (
	"context"
	"database/sql"

	"github.com/FurmanovD/portsvc/pkg/sqldb"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type txCreatorImpl struct {
	db sqldb.SqlDB
}

type txObj struct {
	tx *sql.Tx
}

func NewTxCreator(db sqldb.SqlDB) TxCreator {
	return &txCreatorImpl{
		db: db,
	}
}

func (txc *txCreatorImpl) CreateTransaction(ctx context.Context) (Transaction, error) {
	tx, err := txc.db.Connection().BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &txObj{
		tx: tx,
	}, nil
}

func (r *txObj) Commit() error {
	return r.tx.Commit()
}

func (r *txObj) Rollback() error {
	return r.tx.Rollback()
}

func (r *txObj) Executor() boil.ContextExecutor {
	return r.tx
}
