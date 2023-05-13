package helpers

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

func InstallExtension(extensionName, manifest string, inputs map[string]interface{}, action string) {
	workdir := os.Getenv("WORKDIR")
	repoURL := os.Getenv("EXTENSIONS_REPO")
	branch := os.Getenv("EXTENSIONS_REPO_BRANCH")

	repoPath, repoCloneErr := CloneRepo(repoURL, branch, workdir)
	if repoCloneErr != nil {
		fmt.Printf("Error: %v\n", repoCloneErr)
		return
	}
	logrus.Infof("Cloned repository to: %s\n", repoPath)
	dst := repoPath + "/" + extensionName

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
		installTerraformModule(extensionName, manifestData, inputs2, action, repoPath)
	case "helm":
		installHelmChart(extensionName, manifestData, inputs2, action)
	default:
		logrus.Errorf("Error: Unsupported extension type '%s'.\n", extensionType)
	}
}
