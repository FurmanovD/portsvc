package repository

import (
	"context"

	"github.com/FurmanovD/portsvc/internal/pkg/db/automodel"
)

// PortsRepository contains all functions required to manage Ports objects and their details.
type PortsRepository interface {
	SetPortsInfo(ctx context.Context, ports []*automodel.Port) error
}
