package service

import (
	"context"
)

// PortService is an actual facade of the whole service.
type PortService interface {
	ImportPortsFile(ctx context.Context, portsFilePath string) error
}
