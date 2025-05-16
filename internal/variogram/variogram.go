package variogram

import (
	"fmt"
	"strings"

	"github.com/mmaelicke/go-geostat/internal/types"
)

func NewVariogram(name string, params types.BaseParams) (types.SpatialFunction, error) {
	name = strings.ToLower(name)
	switch name {
	case "spherical":
		return &Spherical{BaseParams: params}, nil
	case "gaussian":
		return &Gaussian{BaseParams: params}, nil
	case "exponential":
		return &Exponential{BaseParams: params}, nil
	case "cubic":
		return &Cubic{BaseParams: params}, nil
	case "matern":
		// Default nu value of 1.5 is a good compromise between smoothness and flexibility
		return &Matern{BaseParams: params, Nu: 1.5}, nil
	default:
		return nil, fmt.Errorf("unknown variogram type: %s", name)
	}
}
