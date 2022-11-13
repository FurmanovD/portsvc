package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/FurmanovD/portsvc/internal/pkg/db/automodel"
	"github.com/FurmanovD/portsvc/pkg/sqldb"
)

type portsRepositoryImpl struct {
	db sqldb.SqlDB
}

func NewPortsRepository(db sqldb.SqlDB) PortsRepository {
	return &portsRepositoryImpl{
		db: db,
	}
}

// ================= interface methods =================================================
func (r *portsRepositoryImpl) SetPortsInfo(ctx context.Context, ports []*automodel.Port) error {
	// NOTE: sqlboiler doesn't support batch insert :(
	var sb strings.Builder
	sb.WriteString("REPLACE INTO ports (`portid`,`name`,`city`,`country`,`province`,`timezone`,`code`,`latitude`,`longitude`) VALUES \n")

	comma := ",\n"
	args := make([]interface{}, 0)
	for i, port := range ports {
		sb.WriteString("(?, ?, ?, ?, ?, ?, ?, ?, ?)")

		if i < len(ports)-1 {
			sb.WriteString(comma)
		}

		var province *string
		if port.Province.Valid {
			province = &ports[i].Province.String
		}
		args = append(args,
			port.Portid, port.Name, port.City, port.Country, province,
			port.Timezone, port.Code, port.Latitude, port.Longitude)

	}
	sb.WriteString(";")

	_, err := r.db.Connection().Exec(sb.String(), args...)
	//_, err := r.db.Connection().Exec(sb.String())
	if err != nil {
		return fmt.Errorf("error executing query '%s' with %+v: %w", sb.String(), ports, err)
	}

	return nil
}
