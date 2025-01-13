package main

import (
	"github.com/CoboGlobal/cobo-mpc-callback-server-v2/cmd/cmd"
	"github.com/CoboGlobal/cobo-mpc-callback-server-v2/internal/config"
	netService "github.com/CoboGlobal/cobo-mpc-callback-server-v2/internal/netservice"
)

const defaultConfigYaml = "configs/callback-server-config.yaml"

func main() {
	cmd.InitDefaultConfig(&config.Config{
		CallbackServer: netService.Config{
			ServiceName:        "callback-server",
			Endpoint:           "0.0.0.0:11020",
			TokenExpireMinutes: 2,
			ClientPubKeyPath:   "configs/tss-node-callback-pub.key",
			ServicePriKeyPath:  "configs/callback-server-pri.pem",
			EnableDebug:        false,
		},
	}, defaultConfigYaml)
	cmd.Execute()
}
