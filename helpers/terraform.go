package helpers

import (
	"fmt"
	"os"
	"os/exec"
)
func mapToString(m map[string]interface{}) []string {
	var keyValuePairs []string

	for key, value := range m {
		keyValuePairs = append(keyValuePairs, fmt.Sprintf("-var=%s=%v", key, value))
	}

	return keyValuePairs
}
func installTerraformModule(extensionName string, manifestData Manifest, dst string, inputs map[string]interface{}, action string) {
	modulePath := fmt.Sprintf("%s/terraform", dst)
	os.Chdir(modulePath)

	cmd := exec.Command("terraform", "init")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error: Failed to initialize Terraform for extension '%s'.\n", extensionName)
		fmt.Println(string(output))
		return
	}
	slice1 := []string{"apply", "-auto-approve"}
	args := append(slice1, mapToString(inputs)...)
	cmd = exec.Command("terraform", args...)
	output, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error: Failed to apply Terraform module for extension '%s'.\n", extensionName)
		fmt.Println(string(output))
		return
	}

	fmt.Printf("Extension '%s' installed successfully.\n", extensionName)
}
