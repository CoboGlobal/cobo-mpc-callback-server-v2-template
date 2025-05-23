package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/CoboGlobal/cobo-mpc-event-server/internal/service"
	"github.com/CoboGlobal/cobo-mpc-event-server/pkg/log"
)

func start() {
	fmt.Println("tss-node-event-server")
	initConfigFile(ConfigFile)

	if CfgInstance == nil {
		log.Fatal("service config empty")
	}

	go trapSignal()

	srv := service.New(CfgInstance)
	srv.Start()
}

func trapSignal() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigs
	exitCode := 128
	if sig == syscall.SIGINT {
		exitCode += int(syscall.SIGINT)
	} else if sig == syscall.SIGTERM {
		exitCode += int(syscall.SIGTERM)
	}
	log.Infof("Received exit signal: %s, code: %v", sig.String(), exitCode)
	os.Exit(exitCode)
}
