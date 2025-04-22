package cmd

import (
	"os"

	"github.com/CoboGlobal/cobo-mpc-event-server/internal/config"
	"github.com/CoboGlobal/cobo-mpc-event-server/pkg/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	CfgInstance       *config.Config
	DefaultConfigYaml string
	ConfigFile        string
)

func InitDefaultConfig(config *config.Config, defaultConfigFile string) {
	CfgInstance = config
	DefaultConfigYaml = defaultConfigFile
}

func InitCmd() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

func AddFlag() {
	rootCmd.PersistentFlags().StringVarP(&ConfigFile, "config", "c", DefaultConfigYaml, "config yaml file path")
}

var rootCmd = &cobra.Command{
	Use:   "tss-node-event-server",
	Short: "TSS Node event server",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		start()
	},
}

func Execute() {
	InitCmd()
	AddFlag()

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

//nolint:nestif
func initConfigFile(config string) {
	if config == "" {
		return
	} else if config == DefaultConfigYaml {
		if _, err := os.Stat(config); err != nil {
			return
		} else {
			log.Infoln("Init from default config file", config)
		}
	} else {
		if _, err := os.Stat(config); err != nil {
			log.Fatalf("Failed to init from config file %v, error: %v", config, err)
		} else {
			log.Infoln("Init from config file", config)
		}
	}
	viper.SetConfigFile(config)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	if viper.ConfigFileUsed() != "" {
		if err := viper.Unmarshal(CfgInstance); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal("config file not used")
	}
}
