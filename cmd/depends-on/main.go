package main

import (
	"context"
	"depends-on/checker"
	"depends-on/resolver"
	"fmt"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {

	ns := resolver.GetNamespace()
	fmt.Printf("Default namespace set to: %s\n", ns)

	fmt.Println("Creating client...")
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	ctx := context.Background()

	// Read annotation of the pod
	annotation, ok := resolver.ReadAnnotation(ctx, clientSet, "depends_on")
	if !ok {
		fmt.Println("No depends-on annotation found.")
		return
	}

	fmt.Printf("Annotation: %s\n", annotation)

	// load dependency
	deps, err := resolver.CheckDependency(annotation)
	if err != nil {
		panic(err.Error())
	}

	filler := resolver.DefaultFillerFactory(ns)

	newCtx, err := checker.LoadAnnotationOptions(ctx, clientSet)
	if err != nil {
		panic(err.Error())
	}
	// Check dependencies
	for _, dep := range deps {
		newDep, err := filler(dep)
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("Checking dependency: %s\n", newDep)
		err = checker.WaitUntilResourceReady(newCtx, clientSet, *newDep.Locator.Namespace, newDep.Locator.Name, newDep.Resource, *newDep.Status)
		if err != nil {
			panic(err.Error())
		}
	}
	fmt.Println("All dependencies are ready.")
}
