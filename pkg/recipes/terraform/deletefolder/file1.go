package main

import (
    "context"
    "fmt"
    "os"
    "path/filepath"

    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/rest"
    "k8s.io/client-go/tools/clientcmd"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)
//Initialize the Kubernetes Client

func getKubernetesClient() (*kubernetes.Clientset, error) {
    // Try in-cluster config
    config, err := rest.InClusterConfig()
    if err != nil {
        // Fallback to kubeconfig
        kubeconfig := filepath.Join(
            homeDir(), ".kube", "config",
        )
        config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
        if err != nil {
            return nil, fmt.Errorf("failed to build kubeconfig: %v", err)
        }
    }

    // Create the clientset
    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        return nil, fmt.Errorf("failed to create clientset: %v", err)
    }

    return clientset, nil
}

func homeDir() string {
    if h := os.Getenv("HOME"); h != "" {
        return h
    }
    return os.Getenv("USERPROFILE") // windows
}
