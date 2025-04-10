package checker

import (
	"context"
	"depends-on/resolver"
	"errors"
	"fmt"
	"k8s.io/client-go/kubernetes"
	"strconv"
	"time"
)

var (
	mapper = map[string]func(context.Context, *kubernetes.Clientset, string, string, string) (bool, error){
		// corev1
		"pod":     CheckPod,
		"service": CheckService,
		// appsv1
		"deployment":  CheckDeployment,
		"statefulset": CheckStatefulSet,
		"replicaset":  CheckReplicaSet,
		"daemonset":   CheckDaemonSet,
	}
)

func CheckResource(ctx context.Context, cs *kubernetes.Clientset, namespace string, name string, resource string, status string) (bool, error) {
	if checkFunc, ok := mapper[resource]; ok {
		return checkFunc(ctx, cs, namespace, name, status)
	}
	return false, InvalidResourceError{resource}
}

func WaitUntilResourceReady(ctx context.Context, cs *kubernetes.Clientset, namespace string, name string, resource string, status string) error {
	tolerance, ok := ctx.Value("tolerance").(int)
	if !ok {
		tolerance = -1
	}
	waitTime, ok := ctx.Value("wait_time").(time.Duration)
	if !ok {
		waitTime = 5 * time.Second
	}

	for {
		avail, err := CheckResource(ctx, cs, namespace, name, resource, status)
		if err != nil {
			var (
				podNotFoundError PodNotFoundError
				notFoundError    NotFoundError
			)
			switch {
			case errors.As(err, &podNotFoundError), errors.As(err, &notFoundError):
				// Not found, possibly not yet created, wait and retry
				fmt.Printf("Resource %s/%s not found, waiting... (tolerance: %d)\n", namespace, name, tolerance)
				if tolerance > 0 {
					tolerance--
				} else if tolerance == -1 {
					// If tolerance is -1, wait indefinitely
				} else {
					return err
				}
			default:
				// Other errors, fail immediately
				return err
			}
		} else {
			if avail {
				break
			}
		}

		// Sleep for a while before checking again
		time.Sleep(waitTime)
	}
	return nil
}

func LoadAnnotationOptions(ctx context.Context, cs *kubernetes.Clientset) (context.Context, error) {
	// Load the annotation options from the resource
	tol, ok := resolver.ReadAnnotation(ctx, cs, "depends-on-tolerance")
	if ok {
		tolerance, err := strconv.Atoi(tol)
		if err != nil {
			return ctx, err
		}
		ctx = context.WithValue(ctx, "tolerance", tolerance)
	}
	waitTime, ok := resolver.ReadAnnotation(ctx, cs, "depends-on-wait-time")
	if ok {
		waitTimeInt, err := strconv.Atoi(waitTime)
		if err != nil {
			return ctx, err
		}
		ctx = context.WithValue(ctx, "wait_time", time.Duration(waitTimeInt)*time.Second)
	}

	return ctx, nil
}
