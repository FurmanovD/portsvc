package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/FurmanovD/portsvc/internal/app/service"
	"github.com/FurmanovD/portsvc/internal/pkg/config"
	"github.com/FurmanovD/portsvc/internal/pkg/db/apidbconvert/v1"
	"github.com/FurmanovD/portsvc/internal/pkg/db/repository"
	"github.com/FurmanovD/portsvc/pkg/commoncfg"
	"github.com/FurmanovD/portsvc/pkg/filewatcher"
	"github.com/FurmanovD/portsvc/pkg/gracefulshutdown"
	"github.com/FurmanovD/portsvc/pkg/sqldb"
	"github.com/sirupsen/logrus"
)

const (
	// program's exit codes.
	errCodeConfigError       = 1
	errCodeDBConnectionError = 2
	errCodeFSWatchingError   = 3
)

// Build information
// The actual information will be stored when 'go build' is called from the Docker file.
var (
	Version   = "local-dev"
	BuildTime = time.Now().Format(time.RFC3339)
	GitCommit = ""

	buildInfo = ""
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})

	buildInfo = fmt.Sprintf(
		"Version: %v BuildTime: %v GitCommit: %v",
		Version,
		BuildTime,
		GitCommit,
	)

	logrus.Info(buildInfo)
}

func main() {

	// Parse flags/config file to populate config.
	cfg, err := config.ParseConfig(os.Args[1:])
	if err != nil {
		fmt.Printf("Configuration load error: %+v", err)
		os.Exit(errCodeConfigError)
	}
	log := initLogging(cfg.LogLevel)
	log.Infof("Logger initialized with LogLevel: %v", cfg.LogLevel)

	log.Infof("configuration loaded: %+v", cfg) // TODO REMOVE ME

	// create a DB connection
	log.Info("Creating a DB connection...")
	dbInstance, err := initDBConnection(sqldb.NewDB(), &cfg.SQLConfig)
	if err != nil {
		log.Errorf("Creating DB connection failed: %v", err)
		os.Exit(errCodeDBConnectionError)
	}

	// create a file watcher instance.
	watcher := filewatcher.New()

	// Create a service instance that will do all required operations to DB, storages etc.
	log.Infof("Creating service instance")
	portsvcService := service.NewService(
		cfg.Service,
		log,
		repository.NewRepository(dbInstance),
		apidbconvert.NewAPIDBConverter(),
	)
	log.Infof("Starting a file watcher to listen for files created in '%s'", cfg.InputDirPath)

	// subscribe for file system events.
	// TODO possibly just pass the watcher interface to the service instance at creation instead.
	err = watcher.Start(
		cfg.InputDirPath,
		func(createdPath string) {
			if err := portsvcService.ImportPortsFile(context.Background(), createdPath); err != nil {
				log.Errorf("Processing file '%s' error: %+v", createdPath, err)
			} else {
				_ = os.Remove(createdPath)
			}
		},
		nil,
		nil,
	)
	if err != nil {
		log.Errorf("Watching filesystem failed: %v", err)
		os.Exit(errCodeFSWatchingError)
	}

	allIsStopped := make(chan struct{})

	// register interception of system signals to gracefully shutdown the service instance.
	gracefulshutdown.InterceptOSSignals(
		gracefulStopAll(log, allIsStopped, watcher, dbInstance),
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGKILL,
	)

	// wait until the graceful shutdown is completed.
	<-allIsStopped
}

// initLogging establishes process logging level.
func initLogging(logLevel string) *logrus.Entry {

	// sets the logging level in app.
	level, err := logrus.ParseLevel(strings.ToLower(logLevel))
	if err != nil {
		level = logrus.InfoLevel
	}
	logrus.SetLevel(level)

	// TODO add host field.
	return logrus.WithField("service", "portsvc")
}

// initDBConnection establishes a connection to a DB.
func initDBConnection(db sqldb.SqlDB, dbConfig *commoncfg.SQLDBConfig) (sqldb.SqlDB, error) {

	err := db.Connect(
		dbConfig.MySQLConnectionString(),
		dbConfig.MaxConnections,
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func gracefulStopAll(
	log *logrus.Entry,
	done chan struct{},
	watcher filewatcher.FileWatcher,
	dbInstance sqldb.SqlDB,
) func() {
	return func() {
		log.Info("Stopping watcher...")
		watcher.Stop()

		log.Info("Closing DB connection...")
		dbInstance.Connection().Close()

		// TODO add any other graceful cleanup calls here.
		done <- struct{}{}
	}
}
