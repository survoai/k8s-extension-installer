# Humalect Extensions Orchestrator

## Introduction

The `k8s-extensions-installer` is a utility tool designed to simplify the installation and management of Kubernetes extensions. This tool allows users to install, update, and delete custom resources, deployments, services, and other Kubernetes objects associated with extensions easily. The `k8s-extensions-installer` supports reading Kubernetes manifest files and applying them to a Kubernetes cluster, with support for templating and variable substitution.

The extensions can be founder here at [Humalect Extensions](https://github.com/Humalect/humalect-extensions)

## Documentation

### Getting Started

1. Clone the repository
```
git clone https://github.com/Humalect/k8s-extension-installer.git
```
2. Build the project
```
go build -o heoctl
```
3. Run the utility
```
./heoctl
```
### Usage

- Install an extension:
```
./heoctl install <extension_name> --inputs <input_variables>

# Example
./heoctl install nginx-k8s --input appname=nginx-deploy,replicas=1
```
- Uninstall an extension:
```
./heoctl uninstall <extension_name> --inputs <input_variables>
```

### Configuration

The `k8s-extensions-installer` requires the following input parameters:

- `<extension-name>`: The name of the extension to be installed or deleted.
- `--inputs`: A map of input variables to be substituted in the manifest files.
