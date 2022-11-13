package service

import (
	"github.com/FurmanovD/portsvc/internal/pkg/db/apidbconvert/v1"
	"github.com/FurmanovD/portsvc/internal/pkg/db/repository"
	"github.com/sirupsen/logrus"
)

type serviceImpl struct {
	cfg       Config
	db        *repository.Repository
	converter apidbconvert.APIDBConverter
	log       *logrus.Entry
}

func NewService(
	cfg Config,
	log *logrus.Entry,
	db *repository.Repository,
	converter apidbconvert.APIDBConverter,
) PortService {
	return &serviceImpl{
		cfg:       cfg,
		log:       log.WithField("subsystem", "port-service"),
		db:        db,
		converter: converter,
	}
}
