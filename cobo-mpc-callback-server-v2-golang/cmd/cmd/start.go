package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/CoboGlobal/cobo-mpc-callback-server-v2/internal/service"
	"github.com/CoboGlobal/cobo-mpc-callback-server-v2/internal/verifier"
	"github.com/CoboGlobal/cobo-mpc-callback-server-v2/pkg/log"
)

func start() {
	fmt.Println("tss-node-callback-server")
	initConfigFile(ConfigFile)

	if CfgInstance == nil {
		log.Fatal("service config empty")
	}

	go trapSignal()
	srv := service.New(CfgInstance, verifier.NewTssVerifier(CfgInstance.AddressWhitelist))
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
