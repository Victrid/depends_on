package checker

import (
	"context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
)

func podReady(pod *v1.Pod) bool {
	for _, condition := range pod.Status.Conditions {
		if condition.Type == v1.PodReady && condition.Status == v1.ConditionTrue {
			return true
		}
	}
	return false
}

func podRunning(pod *v1.Pod) bool {
	return pod.Status.Phase == v1.PodRunning
}

func anyPod(pods []v1.Pod, functor func(*v1.Pod) bool) bool {
	for _, pod := range pods {
		if functor(&pod) {
			return true
		}
	}
	return false
}

func allPod(pods []v1.Pod, functor func(*v1.Pod) bool) bool {
	for _, pod := range pods {
		if !functor(&pod) {
			return false
		}
	}
	return true
}

func CheckPod(ctx context.Context, cs *kubernetes.Clientset, namespace string, name string, status string) (bool, error) {
	resource, err := cs.CoreV1().Pods(namespace).Get(ctx, name, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		return false, NotFoundError{namespace, name}
	}
	if err != nil {
		return false, err
	}
	if resource == nil {
		return false, fmt.Errorf("API error when checking %s/%s", namespace, name)
	}

	switch status {
	case "ready":
		return podReady(resource), nil
	case "running":
		return podRunning(resource), nil
	default:
		return false, InvalidStatusError{name, status}
	}
}

// CheckService checks the status of a service in a Kubernetes cluster.
// This is tricky as services don't maintain status. Available options are:
// - available: exists a pod selected by the service is available.
// - ready: exists a pod selected by the service is ready.
func CheckService(ctx context.Context, cs *kubernetes.Clientset, namespace string, name string, status string) (bool, error) {
	resource, err := cs.CoreV1().Services(namespace).Get(ctx, name, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		return false, NotFoundError{namespace, name}
	}
	if err != nil {
		return false, err
	}
	if resource == nil {
		return false, fmt.Errorf("API error when checking %s/%s", namespace, name)
	}

	// Get the pods selected by the service
	selector := resource.Spec.Selector
	if selector == nil {
		return false, InvalidSelectorError{namespace, name}
	}

	labelSelector := labels.Set(selector).AsSelector().String()

	pods, err := cs.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{LabelSelector: labelSelector})
	if err != nil {
		return false, err
	}
	if len(pods.Items) == 0 {
		return false, PodNotFoundError{namespace, name}
	}

	switch status {
	case "ready", "anyready":
		return anyPod(pods.Items, podReady), nil
	case "running", "anyrunning":
		return anyPod(pods.Items, podRunning), nil
	case "allready":
		return allPod(pods.Items, podReady), nil
	case "allrunning":
		return allPod(pods.Items, podRunning), nil
	default:
		return false, InvalidStatusError{name, status}
	}
}
