package checker

import "fmt"

type (
	NotFoundError struct {
		name      string
		namespace string
	}

	PodNotFoundError struct {
		name      string
		namespace string
	}

	InvalidStatusError struct {
		resource string
		status   string
	}

	InvalidSelectorError struct {
		name      string
		namespace string
	}

	InvalidResourceError struct {
		resource string
	}
)

func (n NotFoundError) Error() string {
	return fmt.Sprintf("resource %s/%s not found", n.namespace, n.name)
}

func (p PodNotFoundError) Error() string {
	return fmt.Sprintf("pod of %s/%s not found", p.namespace, p.name)
}

func (i InvalidStatusError) Error() string {
	return fmt.Sprintf("invalid status (%s) for %s", i.status, i.resource)
}

func (i InvalidSelectorError) Error() string {
	return fmt.Sprintf("Service %s/%s does not have a selector", i.namespace, i.name)
}

func (i InvalidResourceError) Error() string {
	return fmt.Sprintf("invalid resource (%s)", i.resource)
}
