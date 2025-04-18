package main

import (
	"github.com/CoboGlobal/cobo-mpc-event-server/cmd/cmd"
	"github.com/CoboGlobal/cobo-mpc-event-server/internal/config"
	netService "github.com/CoboGlobal/cobo-mpc-event-server/internal/netservice"
)

const defaultConfigYaml = "configs/event-server-config.yaml"

func main() {
	cmd.InitDefaultConfig(&config.Config{
		EventServer: netService.Config{
			ServiceName:      "event-server",
			Endpoint:         "0.0.0.0:11030",
			ClientPubKeyPath: "configs/tss-node-event-pub.key",
			EnableDebug:      false,
		},
	}, defaultConfigYaml)
	cmd.Execute()
}
