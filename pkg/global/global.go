package global

import "github.com/droomlab/drm-coupon/pkg/config"

var Global appGlobal

type appGlobal struct {
	Config config.AppConfig
}

func init() {
	Global.Config = config.Config
}
