package commands

import (
	"github.com/Humalect/k8s-extension-installer/helpers"
	"github.com/spf13/cobra"
)

var UninstallCmd = &cobra.Command{
	Use:   "uninstall [extension-name]",
	Short: "Uninstall a Kubernetes extension",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		extensionName := args[0]
		inputs := parseInputs(inputFlags)
		helpers.InstallExtension(extensionName, manifest, inputs, "uninstall")
	},
}

func init() {
	UninstallCmd.Flags().StringVarP(&inputFlags, "input", "i", "", "Comma-separated list of input key-value pairs")
	UninstallCmd.Flags().StringVarP(&manifest, "manifest", "m", "manifest.yaml", "Path to manifest file")
}
