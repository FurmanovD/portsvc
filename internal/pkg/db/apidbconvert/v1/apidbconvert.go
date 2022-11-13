// package contains interfaces that allows to convert between API and DB object types.
package apidbconvert

import (
	"strconv"

	"github.com/FurmanovD/portsvc/internal/pkg/db/automodel"
	"github.com/FurmanovD/portsvc/internal/pkg/models"
	"github.com/volatiletech/null/v8"
)

type apiDBConverterImpl struct {
}

func NewAPIDBConverter() APIDBConverter {
	return &apiDBConverterImpl{}
}

// ModelToDBPort converts model object to automodel object.
func (c *apiDBConverterImpl) ModelToDBPort(port *models.Port) *automodel.Port {
	if port == nil {
		return nil
	}

	lattitude, longtitude := "", ""
	// TODO what to do in case one of the coordinates is missing?
	if len(port.Coordinates) > 1 {
		lattitude = strconv.FormatFloat(port.Coordinates[0], 'f', 8, 64)
		longtitude = strconv.FormatFloat(port.Coordinates[1], 'f', 8, 64)
	}

	var province null.String
	if port.Province != "" {
		province = null.StringFrom(port.Province)
	}

	return &automodel.Port{
		Portid:    port.PortID,
		Name:      port.Name,
		City:      port.City,
		Country:   port.Country,
		Province:  province,
		Timezone:  port.Timezone,
		Code:      port.Code,
		Latitude:  lattitude,
		Longitude: longtitude,
	}
}

// ModelToDBPorts converts model slice to automodel objects slice.
func (c *apiDBConverterImpl) ModelToDBPorts(items []*models.Port) []*automodel.Port {
	if items == nil {
		return nil
	}

	res := make([]*automodel.Port, len(items))
	for i, item := range items {
		res[i] = c.ModelToDBPort(item)
	}
	return res

}
