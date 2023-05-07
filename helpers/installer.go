package helpers

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type MultiError struct {
	Errors []error
}

func (m *MultiError) Error() string {
	var errorMessages []string
	for _, err := range m.Errors {
		errorMessages = append(errorMessages, err.Error())
	}
	return strings.Join(errorMessages, "; ")
}

func buildInputs(manifestData Manifest, userInput map[string]interface{}) (map[string]interface{}, error) {
	inputs := make(map[string]interface{})
	var multiError MultiError

	for _, input := range manifestData.Inputs {
		// Check if the input value is provided by the user
		if value, ok := userInput[input.Name]; ok {
			// If the input is required, and the user provided a value, use it
			inputs[input.Name] = value
		} else {
			// If the input is required and not provided by the user, add an error
			if input.Required {
				multiError.Errors = append(multiError.Errors, fmt.Errorf("input %s is required, but no value provided", input.Name))
			} else {
				// If the input is not required, use the default value if provided
				inputs[input.Name] = input.Default
			}
		}
	}

	if len(multiError.Errors) > 0 {
		return nil, &multiError
	}
	return inputs, nil
}

func InstallExtension(extensionName, manifest string, inputs map[string]interface{}, action string) {
	fmt.Print(inputs)
	workdir := os.Getenv("WORKDIR")
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	src := "examples/" + extensionName
	dst := workdir + extensionName + "-" + timestamp
	logrus.WithFields(logrus.Fields{
		"src": src,
		"dst": dst,
	}).Info("Directory copying started..............")
	err := CopyDir(src, dst)
	if err != nil {
		logrus.WithError(err).Error("Failed to copy directory\n")
		return
	}
	logrus.Info("Directory copied successfully to ------> " + dst)
	logrus.Info("Reading manifest file..............")
	manifestFile := dst + "/" + manifest
	manifestData, err := readManifest(manifestFile)
	if err != nil {
		panic(err)
	}
	inputs2, err := buildInputs(manifestData, inputs)
	if err != nil {
		logrus.WithError(err).Error("Failed to build inputs.\n")
		return
	}
	logrus.WithFields(logrus.Fields{
		"name": manifestData.Name,
		"type": manifestData.Type,
	}).Info("Reading manifest file completed")

	extensionType := manifestData.Type
	// Print the outputs

	switch extensionType {
	case "k8s":
		kubernetesManifest(extensionName, manifestData, dst, inputs2, action)
	case "terraform":
		installTerraformModule(extensionName, manifestData, inputs2, action)
	case "helm":
		installHelmChart(extensionName, manifestData, inputs2, action)
	default:
		logrus.Errorf("Error: Unsupported extension type '%s'.\n", extensionType)
	}
}
