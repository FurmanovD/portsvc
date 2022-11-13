package repository

import (
	"github.com/FurmanovD/portsvc/pkg/sqldb"
)

// Repository implements the Repository interface.
type Repository struct {
	TxCreator TxCreator
	Ports     PortsRepository
}

func NewRepository(db sqldb.SqlDB) *Repository {
	return &Repository{
		TxCreator: NewTxCreator(db),
		Ports:     NewPortsRepository(db),
	}
}
