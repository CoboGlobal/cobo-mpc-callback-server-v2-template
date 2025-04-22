package config

import (
	netService "github.com/CoboGlobal/cobo-mpc-event-server/internal/netservice"
)

type Config struct {
	EventServer netService.Config `mapstructure:"event_server"`
}
