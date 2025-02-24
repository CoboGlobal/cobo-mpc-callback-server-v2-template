package config

import (
	netService "github.com/CoboGlobal/cobo-mpc-callback-server-v2/internal/netservice"
)

type Config struct {
	CallbackServer   netService.Config `mapstructure:"callback_server"`
	AddressWhitelist []string          `mapstructure:"address_whitelist"`
}
