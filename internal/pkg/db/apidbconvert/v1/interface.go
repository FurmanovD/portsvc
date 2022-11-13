package apidbconvert

import (
	"github.com/FurmanovD/portsvc/internal/pkg/db/automodel"
	"github.com/FurmanovD/portsvc/internal/pkg/models"
)

type APIDBConverter interface {
	// these functions convert a DB structure(s) to API object(s) and vice versa
	ModelToDBPort(port *models.Port) *automodel.Port
	ModelToDBPorts(ports []*models.Port) []*automodel.Port
}
