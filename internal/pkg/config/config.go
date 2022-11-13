package config

import (
	"flag"

	"github.com/peterbourgon/ff"
	"github.com/peterbourgon/ff/fftoml"
	"github.com/sirupsen/logrus"

	"github.com/FurmanovD/portsvc/internal/app/service"
	"github.com/FurmanovD/portsvc/pkg/commoncfg"
)

const (
	// ConfigPrefix used in reading the base constants from ini file
	ConfigPrefix = "portsvc."

	// Default INI file to use
	DefaultConfigIni = "/app/config.ini"

	// default log level to use
	DefaultLogLevel = "info"

	// default DB maxConnections to use
	DefaultDBMaxConnections = 10
)

// Config stuct for the Station API
type Config struct {
	LogLevel     string
	InputDirPath string

	SQLConfig commoncfg.SQLDBConfig
	Service   service.Config
}

// ParseConfig parses the configuration file
func ParseConfig(args []string) (*Config, error) {
	fs := flag.NewFlagSet("portsvc", flag.ContinueOnError)

	cfg := &Config{}

	// initialize parameters:
	fs.StringVar(&cfg.InputDirPath, ConfigPrefix+"INPUT_DIRECTORY_PATH", "/in", "file watcher path()")

	fs.StringVar(
		&cfg.LogLevel,
		"loglevel",
		"info",
		"Use this option to set the log level for the application. "+
			"Possible Values: trace, debug, info, warn, error, fatal, panic, nolog",
	)
	// end of command-line parameters

	initDBConfig(fs, &cfg.SQLConfig)

	configFilePath := fs.String("config", DefaultConfigIni, "a config file name")
	filePath := DefaultConfigIni
	if configFilePath != nil {
		filePath = *configFilePath
	}
	logrus.Infof("config file to be loaded:%s", filePath)

	err := ff.Parse(fs, args,
		ff.WithEnvVarNoPrefix(),
		ff.WithConfigFileFlag("config"),
		ff.WithConfigFileParser(fftoml.Parser),
		ff.WithIgnoreUndefined(true),
	)

	return cfg, err
}

// initDBConfig parses the DB configuration parameters
func initDBConfig(fs *flag.FlagSet, cfg *commoncfg.SQLDBConfig) {
	//  MySQL config
	fs.StringVar(&cfg.Host, ConfigPrefix+"DB_HOST", "", "host name of mysql db instance")
	fs.IntVar(&cfg.Port, ConfigPrefix+"DB_PORT", commoncfg.DefaultMySQLPort, "port number name of mysql db instance")
	fs.StringVar(&cfg.Database, ConfigPrefix+"DB_NAME", "", "write instance database name")
	fs.StringVar(&cfg.User, ConfigPrefix+"DB_USER", "", "mysql user name for db")
	fs.StringVar(&cfg.Password, ConfigPrefix+"DB_PASS", "", "password for mysql user for db")
	fs.IntVar(&cfg.MaxConnections, ConfigPrefix+"DB_MAX_CONNECTIONS", DefaultDBMaxConnections, "max connections to a DB")
}
