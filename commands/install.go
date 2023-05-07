package commands

import (
	"os"
	"strings"

	"github.com/Humalect/k8s-extension-installer/helpers"
	"github.com/spf13/cobra"
)

var inputFlags string
var manifest string

func parseInputs(inputFlags string) map[string]interface{} {
	inputs := make(map[string]interface{})

	if inputFlags != "" {
		pairs := strings.Split(inputFlags, ",")
		for _, pair := range pairs {
			kv := strings.Split(pair, "=")
			if len(kv) == 2 {
				inputs[kv[0]] = kv[1]
			}
		}
	}

	for _, envVar := range os.Environ() {
		kv := strings.Split(envVar, "=")
		if len(kv) == 2 && strings.HasPrefix(kv[0], "EXT_") {
			inputs[kv[0][4:]] = kv[1]
		}
	}

	return inputs
}

var InstallCmd = &cobra.Command{
	Use:   "install [extension-name]",
	Short: "Install a Kubernetes extension",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		extensionName := args[0]
		inputs := parseInputs(inputFlags)
		helpers.InstallExtension(extensionName, manifest, inputs, "install")
	},
}

func init() {
	InstallCmd.Flags().StringVarP(&inputFlags, "input", "i", "", "Comma-separated list of input key-value pairs")
	InstallCmd.Flags().StringVarP(&manifest, "manifest", "m", "manifest.yaml", "Path to manifest file")
}
