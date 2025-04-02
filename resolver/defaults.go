package resolver

import "fmt"

var (
	resourceDefaultMap = map[string]string{
		"pod":         "ready",
		"service":     "ready",
		"deployment":  "available",
		"statefulset": "available",
		"replicaset":  "available",
		"daemonset":   "available",
	}
)

func DefaultFillerFactory(namespace string) func(dependency *Dependency) (*Dependency, error) {
	return func(dependency *Dependency) (*Dependency, error) {
		if dependency.Locator.Namespace == nil {
			dependency.Locator.Namespace = &namespace
		}
		if dependency.Status == nil {
			if status, ok := resourceDefaultMap[dependency.Resource]; ok {
				dependency.Status = &status
			} else {
				return nil, fmt.Errorf("resource %s not found in default map", dependency.Resource)
			}
		}
		return dependency, nil
	}
}
