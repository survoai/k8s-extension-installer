package helpers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/template"

	"github.com/sirupsen/logrus"
	gopkgYaml "gopkg.in/yaml.v3"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
)

func getKubeConfig() (*rest.Config, error) {
	kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func kubernetesManifest(extensionName string, manifestData Manifest, destination string, inputs map[string]interface{}, action string) error {
	config, err := getKubeConfig()
	if err != nil {
		logrus.WithError(err).Error("Error: Failed to load kubeconfig:\n")
		return err
	}

	discoveryClient, err := getDiscoveryClient(config)
	if err != nil {
		return err
	}

	dynamicClient, err := getDynamicClient(config)
	if err != nil {
		return err
	}

	rm, err := getRESTMapper(discoveryClient)
	if err != nil {
		return err
	}

	err = processManifestFiles(destination, inputs, dynamicClient, rm, action)

	if err != nil {
		return err
	}

	return nil
}

func getDiscoveryClient(config *rest.Config) (*discovery.DiscoveryClient, error) {
	discoveryClient, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		logrus.WithError(err).Error("Error: Failed to create discovery client:\n")
		return nil, err
	}
	return discoveryClient, nil
}

func getDynamicClient(config *rest.Config) (dynamic.Interface, error) {
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		logrus.WithError(err).Error("Error: Failed to create dynamic client:\n")
		return nil, err
	}
	return dynamicClient, nil
}

func readManifestsFromFile(path string, inputs map[string]interface{}) ([]string, error) {
	tmpl, err := template.ParseFiles(path)
	if err != nil {
		logrus.WithError(err).Error("Error: Failed to parse Kubernetes manifest template:\n")
		return nil, err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, inputs)
	if err != nil {
		logrus.WithError(err).Error("Error: Failed to substitute values in Kubernetes manifest template:\n")
		return nil, err
	}

	decoder := gopkgYaml.NewDecoder(&buf)
	var manifestDocs []string

	for {
		var doc gopkgYaml.Node
		err := decoder.Decode(&doc)
		if err == io.EOF {
			break
		}
		if err != nil {
			// Log the error, but return an empty slice and a nil error
			logrus.WithError(err).Error("Error: Failed to decode YAML document:\n")
			return []string{}, nil
		}

		var manifestBuffer bytes.Buffer
		encoder := gopkgYaml.NewEncoder(&manifestBuffer)
		err = encoder.Encode(&doc)
		if err != nil {
			logrus.WithError(err).Error("Error: Failed to encode YAML document:\n")
			return nil, err
		}
		encoder.Close()

		manifestDocs = append(manifestDocs, manifestBuffer.String())
	}

	return manifestDocs, nil
}

func processManifestFiles(destination string, inputs map[string]interface{}, dynamicClient dynamic.Interface, rm meta.RESTMapper, action string) error {
	var resourceStatus sync.Map

	err := filepath.Walk(filepath.Join(destination, "k8s"), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		manifests, err := readManifestsFromFile(path, inputs)
		if err != nil {
			logrus.WithError(err).Error("Error: Failed to process manifest file:\n")
			return nil // Return nil instead of err to continue processing other files
		}

		for _, manifest := range manifests {
			// Add this line to check if the manifest slice is empty
			if len(manifest) == 0 {
				continue
			}

			if action == "delete" || action == "uninstall" {
				err = deleteResource(manifest, path, dynamicClient, rm)
			} else {
				err = installResource(manifest, path, dynamicClient, rm)
			}

			obj, gvk, _ := decodeManifest(manifest)
			resourceKey := strings.Join([]string{path, obj.GetName(), gvk.Kind}, "-")
			if err != nil {
				resourceStatus.Store(resourceKey, ResourceStatus{
					Success: false,
					Message: fmt.Sprintf("Failed %s : %v", action,
						err),
					Resource: obj.GetName(),
					Kind:     gvk.Kind,
					Path:     path,
				})
			} else {
				resourceStatus.Store(resourceKey, ResourceStatus{
					Success:  true,
					Message:  fmt.Sprintf("Success %s : Resource %s of kind %s", action, obj.GetName(), gvk.Kind),
					Resource: obj.GetName(),
					Kind:     gvk.Kind,
					Path:     path,
				})
			}

		}

		return nil // Return nil to continue processing other files
	})

	if err != nil {
		logrus.WithError(err).Error("Error: Failed to process manifest files:\n")
		return err
	}

	logrus.Info("========================= Resource installation summary ========================= \n")
	resourceStatus.Range(func(key, value interface{}) bool {
		resourceStat := value.(ResourceStatus)

		if resourceStat.Success {
			logrus.WithFields(logrus.Fields{
				"Name": resourceStat.Resource,
				"Kind": resourceStat.Kind,
				"Path": resourceStat.Path,
			}).Info(resourceStat.Message)
		} else {
			logrus.WithFields(logrus.Fields{
				"Name": resourceStat.Resource,
				"Kind": resourceStat.Kind,
				"Path": resourceStat.Path,
			}).Error(resourceStat.Message)
		}

		return true
	})

	return nil
}

func getRESTMapper(discoveryClient *discovery.DiscoveryClient) (meta.RESTMapper, error) {
	groupResources, err := restmapper.GetAPIGroupResources(discoveryClient)
	if err != nil {
		logrus.WithError(err).Error("Error: Failed to get API group resources:\n")
		return nil, err
	}

	rm := restmapper.NewDiscoveryRESTMapper(groupResources)
	return rm, nil
}

func getResourceClients(manifest string, dynamicClient dynamic.Interface, rm meta.RESTMapper) (dynamic.ResourceInterface, *unstructured.Unstructured, *schema.GroupVersionKind, error) {
	if strings.TrimSpace(manifest) == "" {
		return nil, nil, nil, nil
	}

	obj, gvk, err := decodeManifest(manifest)
	if err != nil {
		return nil, nil, nil, err
	}

	gvr, err := getGroupVersionResource(&gvk, rm)
	if err != nil {
		return nil, nil, nil, err
	}

	namespaced := gvk.GroupKind().String() != "Namespace"
	var dynamicResourceClient dynamic.ResourceInterface

	if namespaced {
		namespace := obj.GetNamespace()
		if namespace == "" {
			namespace = "default"
		}
		dynamicResourceClient = dynamicClient.Resource(gvr).Namespace(namespace)
	} else {
		dynamicResourceClient = dynamicClient.Resource(gvr)
	}

	return dynamicResourceClient, obj, &gvk, nil
}

func installResource(manifest string, path string, dynamicClient dynamic.Interface, rm meta.RESTMapper) error {
	dynamicResourceClient, obj, _, err := getResourceClients(manifest, dynamicClient, rm)
	if err != nil {
		return err
	}

	if dynamicResourceClient == nil {
		return nil
	}

	_, err = dynamicResourceClient.Create(context.Background(), obj, metav1.CreateOptions{})
	return err
}

func deleteResource(manifest string, path string, dynamicClient dynamic.Interface, rm meta.RESTMapper) error {
	dynamicResourceClient, obj, _, err := getResourceClients(manifest, dynamicClient, rm)
	if err != nil {
		return err
	}

	if dynamicResourceClient == nil {
		return nil
	}

	deletePolicy := metav1.DeletePropagationForeground
	deleteOptions := metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}
	err = dynamicResourceClient.Delete(context.Background(), obj.GetName(), deleteOptions)
	return err
}

func decodeManifest(manifest string) (*unstructured.Unstructured, schema.GroupVersionKind, error) {
	decode := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	obj := &unstructured.Unstructured{}

	_, gvk, err := decode.Decode([]byte(manifest), nil, obj)
	if err != nil {
		return nil, schema.GroupVersionKind{}, err
	}

	return obj, *gvk, nil
}

func getGroupVersionResource(gvk *schema.GroupVersionKind, rm meta.RESTMapper) (schema.GroupVersionResource, error) {
	gvr, err := rm.ResourceFor(schema.GroupVersionResource{Group: gvk.Group, Version: gvk.Version, Resource: gvk.Kind})
	if err != nil {
		return schema.GroupVersionResource{}, err
	}
	return gvr, nil
}
