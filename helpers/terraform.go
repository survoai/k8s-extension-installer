package helpers

import (
	"fmt"
	"os"
	"os/exec"
)

func installTerraformModule(extensionName string, manifestData Manifest, inputs map[string]interface{}, action string) {
	modulePath := fmt.Sprintf("extensions/%s/terraform-module", extensionName)
	os.Chdir(modulePath)

	cmd := exec.Command("terraform", "init")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error: Failed to initialize Terraform for extension '%s'.\n", extensionName)
		fmt.Println(string(output))
		return
	}

	cmd = exec.Command("terraform", "apply", "-auto-approve")
	output, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error: Failed to apply Terraform module for extension '%s'.\n", extensionName)
		fmt.Println(string(output))
		return
	}

	fmt.Printf("Extension '%s' installed successfully.\n", extensionName)
}
