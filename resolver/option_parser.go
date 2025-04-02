package resolver

import (
	"fmt"
	"strings"
)

func CheckDependency(input string) ([]*Dependency, error) {
	input = strings.TrimSpace(input)

	if strings.HasPrefix(input, "ConfigMap=") {
		locatorStr := strings.TrimPrefix(input, "ConfigMap=")
		locator, err := parseLocator(locatorStr)
		if err != nil {
			return nil, fmt.Errorf("invalid ConfigMap locator: %v", err)
		}

		return ReadConfigMap(*locator)
	}

	return parseDependsOnString(input)
}

func parseDependsOnString(input string) ([]*Dependency, error) {
	var dependencies []*Dependency

	parts := strings.Split(input, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		dependency, err := parseDependency(part)
		if err != nil {
			return nil, fmt.Errorf("invalid dependency '%s': %v", part, err)
		}
		dependencies = append(dependencies, dependency)
	}

	return dependencies, nil
}

func parseDependency(input string) (*Dependency, error) {
	parts := strings.SplitN(input, ":", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("expected 'resource:locator' format")
	}

	resource := strings.TrimSpace(parts[0])

	locatorAndStatus := strings.TrimSpace(parts[1])
	var status *string
	if strings.Contains(locatorAndStatus, "?") {
		statusParts := strings.SplitN(locatorAndStatus, "?", 2)
		locatorStr := strings.TrimSpace(statusParts[0])
		statusStr := strings.TrimSpace(statusParts[1])

		status = &statusStr
		locatorAndStatus = locatorStr
	}

	locator, err := parseLocator(locatorAndStatus)
	if err != nil {
		return nil, fmt.Errorf("invalid locator: %v", err)
	}

	return &Dependency{
		Resource: resource,
		Locator:  *locator,
		Status:   status,
	}, nil
}

func parseLocator(input string) (*Locator, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil, fmt.Errorf("empty locator")
	}

	if strings.Contains(input, "/") {
		parts := strings.SplitN(input, "/", 2)
		namespace := strings.TrimSpace(parts[0])
		name := strings.TrimSpace(parts[1])

		if namespace == "" || name == "" {
			return nil, fmt.Errorf("invalid namespace/name format")
		}

		return &Locator{
			Namespace: &namespace,
			Name:      name,
		}, nil
	}

	// 只有name的情况
	return &Locator{
		Namespace: nil,
		Name:      input,
	}, nil
}
