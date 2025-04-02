package resolver

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type YamlDependency struct {
	Resource  *string `yaml:"resource,omitempty"`
	Name      *string `yaml:"name,omitempty"`
	Namespace *string `yaml:"namespace,omitempty"`
	Status    *string `yaml:"status,omitempty"`
	Raw       *string `yaml:"raw,omitempty"`
}

// Read configMap from k8s and parse it

func parseConfigMapDependency(input string) ([]*Dependency, error) {
	// Parse yamlContent
	var yamlDeps []YamlDependency
	err := yaml.Unmarshal([]byte(input), &yamlDeps)
	if err != nil {
		return nil, err
	}

	// Convert yamlDeps to Dependencies
	var deps []*Dependency

	for _, yamlDep := range yamlDeps {
		if yamlDep.Raw != nil {
			// Raw dependency
			dep, err := parseDependency(*yamlDep.Raw)
			if err != nil {
				return nil, err
			}
			deps = append(deps, dep)
		} else {
			// Normal dependency
			if yamlDep.Resource == nil || yamlDep.Name == nil {
				return nil, fmt.Errorf("resource and Name are required for a dependency")
			}
			dep := &Dependency{
				Resource: *yamlDep.Resource,
				Locator: Locator{
					Namespace: yamlDep.Namespace,
					Name:      *yamlDep.Name,
				},
				Status: yamlDep.Status,
			}
			deps = append(deps, dep)
		}
	}

	return deps, nil
}

func ReadConfigMap(locator Locator) ([]*Dependency, error) {
	// Read configMap from k8s
	name := locator.Name
	namespace := ""
	if locator.Namespace == nil {
		namespace = GetNamespace()
	} else {
		namespace = *locator.Namespace
	}

	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	ctx := context.Background()

	configMap, err := clientSet.CoreV1().ConfigMaps(namespace).Get(ctx, name, metav1.GetOptions{})

	if err != nil {
		return nil, err
	}

	configMapData := configMap.Data

	yamlContent, ok := configMapData["depends_on"]
	if !ok {
		return nil, fmt.Errorf("ConfigMap %s/%s does not contain 'depends_on' key", namespace, name)
	}

	return parseConfigMapDependency(yamlContent)
}
