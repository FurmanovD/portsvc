package commoncfg

import (
	"strconv"
)

const (
	// DefaultMySQLPort is default port to use.
	DefaultMySQLPort = 3306
)

// MySQLConnectionString constructs a MySQL connection string from configuration parameters.
func (c *SQLDBConfig) MySQLConnectionString() string {
	// port, 3306 by default.
	p := c.Port
	if p <= 0 {
		p = DefaultMySQLPort
	}
	// host, 127.0.0.1 by default.
	h := c.Host
	if h == "" {
		h = "127.0.0.1"
	}
	return c.User + ":" + c.Password + "@tcp(" + h + ":" + strconv.Itoa(p) + ")/" + c.Database
}
