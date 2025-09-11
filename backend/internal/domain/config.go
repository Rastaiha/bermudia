package domain

import (
	"github.com/Rastaiha/bermudia/internal/config"
)

func ApplyConfig(cfg config.Config) {
	if cfg.DevMode {
		initialKeyCount = 5
		playerOpenOffersLimit = 10
	}
}
