package checker

import (
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func CheckDeployment(ctx context.Context, cs *kubernetes.Clientset, namespace string, name string, status string) (bool, error) {
	resource, err := cs.AppsV1().Deployments(namespace).Get(ctx, name, metav1.GetOptions{})
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
	case "allready":
		return resource.Status.ReadyReplicas == *resource.Spec.Replicas, nil
	case "allavailable":
		return resource.Status.ReadyReplicas == *resource.Spec.Replicas, nil
	case "anyready", "ready":
		return resource.Status.ReadyReplicas > 0, nil
	case "anyavailable", "available":
		return resource.Status.AvailableReplicas > 0, nil
	default:
		return false, InvalidStatusError{name, status}
	}
}

func CheckStatefulSet(ctx context.Context, cs *kubernetes.Clientset, namespace string, name string, status string) (bool, error) {
	// Check if the deployment exists
	resource, err := cs.AppsV1().StatefulSets(namespace).Get(ctx, name, metav1.GetOptions{})
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
	case "allready":
		return resource.Status.ReadyReplicas == *resource.Spec.Replicas, nil
	case "allavailable":
		return resource.Status.ReadyReplicas == *resource.Spec.Replicas, nil
	case "anyready", "ready":
		return resource.Status.ReadyReplicas > 0, nil
	case "anyavailable", "available":
		return resource.Status.AvailableReplicas > 0, nil
	default:
		return false, InvalidStatusError{name, status}
	}
}

func CheckReplicaSet(ctx context.Context, cs *kubernetes.Clientset, namespace string, name string, status string) (bool, error) {
	// Check if the deployment exists
	resource, err := cs.AppsV1().ReplicaSets(namespace).Get(ctx, name, metav1.GetOptions{})
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
	case "allready":
		return resource.Status.ReadyReplicas == *resource.Spec.Replicas, nil
	case "allavailable":
		return resource.Status.ReadyReplicas == *resource.Spec.Replicas, nil
	case "anyready", "ready":
		return resource.Status.ReadyReplicas > 0, nil
	case "anyavailable", "available":
		return resource.Status.AvailableReplicas > 0, nil
	default:
		return false, InvalidStatusError{name, status}
	}
}

func CheckDaemonSet(ctx context.Context, cs *kubernetes.Clientset, namespace string, name string, status string) (bool, error) {
	// Check if the deployment exists
	resource, err := cs.AppsV1().DaemonSets(namespace).Get(ctx, name, metav1.GetOptions{})
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
	case "anyready", "ready":
		// Check if any pod is ready
		return resource.Status.NumberReady > 0, nil
	case "anyavailable", "available":
		return resource.Status.NumberAvailable > 0, nil
	default:
		return false, InvalidStatusError{name, status}
	}
}
