package helpers

import (
	"io"
	"os"

	"gopkg.in/yaml.v2"
)

func readManifest(manifestFile string) (Manifest, error) {
	// Code to read manifestFile
	file, err := os.Open(manifestFile)
	if err != nil {
		return Manifest{}, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return Manifest{}, err
	}

	// Unmarshal the manifest data into a Manifest struct
	var manifestData Manifest
	err = yaml.Unmarshal(data, &manifestData)
	if err != nil {
		return Manifest{}, err
	}

	return manifestData, nil
}
