package resolver

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"os"
)

const (
	annotationPrefix = "victrid.dev/"
)

func GetNamespace() string {
	nsb, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	if err != nil {
		panic(err.Error())
	}
	ns := string(nsb)
	return ns
}

func ReadAnnotation(ctx context.Context, cs *kubernetes.Clientset, name string) (string, bool) {
	// Our pod name from env HOSTNAME
	podName := os.Getenv("HOSTNAME")
	if podName == "" {
		panic("HOSTNAME environment variable not set")
	}

	// Read the annotations through the Kubernetes API
	pod, err := cs.CoreV1().Pods(GetNamespace()).Get(ctx, podName, metav1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}

	// Check if the pod has annotations
	annotation, ok := pod.Annotations[annotationPrefix+name]
	return annotation, ok
}
