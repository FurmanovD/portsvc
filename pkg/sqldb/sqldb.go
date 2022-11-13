package sqldb

import (
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

const (
	DriverMySQL = "mysql"

	connStringParams = "parseTime=true"
)

// SqlDB interface for db.
type SqlDB interface {
	Connect(connStr string, maxConn int) error
	Connection() *sqlx.DB
	ConnectionString() string
}

// db implements the SqlDB interface.
type db struct {
	connStr string
	conn    *sqlx.DB
}

// NewDB returns the new database object.
func NewDB() SqlDB {
	return &db{}
}

// Connect connects to a db using a connection string.
func (d *db) Connect(connStr string, maxConn int) error {
	if err := d.createConnection(connStr, connStringParams, maxConn); err != nil {
		return err
	}
	return nil
}

// Connection returns the current connection.
func (d *db) Connection() *sqlx.DB {
	return d.conn
}

// ConnectionString returns the current connection string.
func (d *db) ConnectionString() string {
	return d.connStr
}

func (d *db) createConnection(connStr string, params string, maxConn int) error {
	fullConnStr := connStr
	if strings.Contains(fullConnStr, "?") {
		if fullConnStr[len(fullConnStr)-1] != '?' {
			fullConnStr += "&"
		}
		fullConnStr += params
	} else {
		fullConnStr += "?" + params
	}

	conn, err := sqlx.Connect(DriverMySQL, fullConnStr)
	if err != nil {
		return err
	}

	d.connStr = fullConnStr

	conn.SetMaxOpenConns(maxConn)
	conn.SetMaxIdleConns(0)

	d.conn = conn

	return nil
}
