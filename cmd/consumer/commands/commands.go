package commands

import (
	"github.com/droomlab/drm-coupon/internal/app/dependency"
	"github.com/droomlab/drm-coupon/pkg/drmrmq"
)

type commandFunc func(*dependency.Dependency) drmrmq.Handler

// Commands contains all delivery processing fuction
// with string key.
type Commands map[string]commandFunc

// Return all available string command map.
func GetCommands() Commands {
	return Commands{
		"HandleDRMTesting": HandleDRMTesting,
	}
}
