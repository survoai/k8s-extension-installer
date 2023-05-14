package helpers

import (
	"fmt"
	"os/exec"
)

func installHelmChart(extensionName string, manifestData Manifest, dst string, inputs map[string]interface{}, action string) {
	chartPath := fmt.Sprintf("%s/helm", dst)
	cmd := exec.Command("helm", "install", extensionName, chartPath)

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error: Failed to install Helm chart for extension '%s'.\n", extensionName)
		fmt.Println(string(output))
		return
	}

	fmt.Printf("Extension '%s' installed successfully.\n", extensionName)
}
