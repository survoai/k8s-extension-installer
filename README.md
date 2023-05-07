# k8s-extensions-installer

## Introduction

The `k8s-extensions-installer` is a utility tool designed to simplify the installation and management of Kubernetes extensions. This tool allows users to install, update, and delete custom resources, deployments, services, and other Kubernetes objects associated with extensions easily. The `k8s-extensions-installer` supports reading Kubernetes manifest files and applying them to a Kubernetes cluster, with support for templating and variable substitution.

## Documentation

### Getting Started

1. Clone the repository
```
git clone https://github.com/example/k8s-extensions-installer.git
```
2. Build the project
```
./k8s-extensions-installer
go build -o k8s-extensions-installer
```
3. Run the utility
```
./k8s-extensions-installer
```
### Usage

- Install an extension:
```
./k8s-extensions-installer install --extension-name <extension_name> --manifest-data <manifest_data_path> --destination <destination_path> --inputs <input_variables>
```
- Delete an extension:
```
./k8s-extensions-installer delete --extension-name <extension_name> --manifest-data <manifest_data_path> --destination <destination_path> --inputs <input_variables>
```

### Configuration

The `k8s-extensions-installer` requires the following input parameters:

- `--extension-name`: The name of the extension to be installed or deleted.
- `--manifest-data`: The path to the manifest data file containing Kubernetes resources for the extension.
- `--destination`: The destination path where the extension will be installed or deleted.
- `--inputs`: A map of input variables to be substituted in the manifest files.
