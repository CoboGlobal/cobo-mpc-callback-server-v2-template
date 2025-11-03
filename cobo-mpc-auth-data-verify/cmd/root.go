package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/CoboGlobal/cobo-mpc-callback-server-v2-template/cobo-mpc-auth-data-verify/validator"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	BizDataFile  string
	TemplateFile string
	MessageFile  string
)

func InitCmd() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

func AddFlag() {
	rootCmd.Flags().StringVar(&BizDataFile, "biz-data-file", "", "input biz data file")
	if err := rootCmd.MarkFlagRequired("biz-data-file"); err != nil {
		log.Fatal(err)
	}
	rootCmd.Flags().StringVar(&TemplateFile, "template-file", "", "input template file")
	if err := rootCmd.MarkFlagRequired("template-file"); err != nil {
		log.Fatal(err)
	}
	rootCmd.Flags().StringVar(&MessageFile, "message-file", "", "output message file")
}

var rootCmd = &cobra.Command{
	Use:   "cobo-auth-data-validator",
	Short: "Validate auth data by biz data and template",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		buildStatementMessage()
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

func buildStatementMessage() {
	fmt.Println("Build statement message by biz data and template")
	fmt.Printf("Biz data file: %v\n", BizDataFile)
	fmt.Printf("Template file: %v\n", TemplateFile)

	BizDataContent, err := os.ReadFile(filepath.Clean(BizDataFile))
	if err != nil {
		log.Fatalf("Read biz data file %v failed: %v", BizDataFile, err)
	}

	TemplateContent, err := os.ReadFile(filepath.Clean(TemplateFile))
	if err != nil {
		log.Fatalf("Read template file %v failed: %v", TemplateFile, err)
	}

	authData := &validator.AuthData{
		Template: string(TemplateContent),
		BizData:  string(BizDataContent),
	}

	authValidator := validator.NewAuthValidator(authData)
	buildMsg, err := authValidator.BuildStatementMessage()
	if err != nil {
		log.Fatal("Failed to build statement message: ", err)
	}

	if MessageFile != "" {
		err := os.MkdirAll(filepath.Dir(MessageFile), 0755)
		if err != nil {
			log.Fatalf("Create message file directory %v failed: %v", MessageFile, err)
		}
		err = os.WriteFile(filepath.Clean(MessageFile), []byte(buildMsg), 0644)
		if err != nil {
			log.Fatalf("Write message file %v failed: %v", MessageFile, err)
		}
		fmt.Printf("Write statement message to file %v\n", MessageFile)
	}

	fmt.Println("Statement message:")
	fmt.Println(buildMsg)
}
