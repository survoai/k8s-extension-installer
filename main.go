package main

import (
	"log"
	"os"

	"github.com/Humalect/k8s-extension-installer/commands"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var rootCmd = &cobra.Command{
		Use:   "heoctl",
		Short: "heoctl is a CLI tool to install and uninstall Kubernetes extensions",
	}

	rootCmd.AddCommand(commands.InstallCmd)
	rootCmd.AddCommand(commands.UninstallCmd)

	if err := rootCmd.Execute(); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
}
