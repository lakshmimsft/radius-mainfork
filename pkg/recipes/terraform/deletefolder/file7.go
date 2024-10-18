// complete eg.. integrating into Deploy()
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "time"

    corev1 "k8s.io/api/core/v1"
    batchv1 "k8s.io/api/batch/v1"
    "k8s.io/apimachinery/pkg/api/errors"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/api/resource"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/rest"
    "k8s.io/client-go/tools/clientcmd"
)

type Options struct {
    WorkingDir              string
    RecipeName              string
    Namespace               string
    CloudProviderCredentials map[string]string
    // Add other fields as necessary.
}

type executor struct {
    clientset *kubernetes.Clientset
    // Add other fields if necessary.
}

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

func (e *executor) Deploy(ctx context.Context, options Options) error {
    // Initialize Kubernetes client
    clientset, err := getKubernetesClient()
    if err != nil {
        return fmt.Errorf("failed to get Kubernetes client: %v", err)
    }

    // Define PV and PVC
    pv := definePersistentVolume()
    pvc := definePersistentVolumeClaim(options.Namespace)

    // Create PV if it doesn't exist
    if err := createPersistentVolume(ctx, clientset, pv); err != nil {
        return err
    }

    // Create PVC if it doesn't exist
    if err := createPersistentVolumeClaim(ctx, clientset, pvc); err != nil {
        return err
    }

    // Proceed with generating config, creating ConfigMap, Secret, and Job
    if err := e.generateConfig(ctx, options); err != nil {
        return fmt.Errorf("failed to generate config: %v", err)
    }

    configMapName := fmt.Sprintf("terraform-config-%s", options.RecipeName)
    if err := e.createConfigMap(ctx, configMapName, options); err != nil {
        return fmt.Errorf("failed to create ConfigMap: %v", err)
    }

    envVars, err := e.getEnvironmentVariables(ctx, options)
    if err != nil {
        return fmt.Errorf("failed to get environment variables: %v", err)
    }

    secretName := fmt.Sprintf("terraform-secret-%s", options.RecipeName)
    if err := e.createSecret(ctx, secretName, envVars, options.Namespace); err != nil {
        return fmt.Errorf("failed to create Secret: %v", err)
    }

    // Create Kubernetes Job
    jobName := fmt.Sprintf("terraform-apply-job-%s", options.RecipeName)
    if err := e.createTerraformJob(ctx, jobName, configMapName, secretName, "apply", options); err != nil {
        return fmt.Errorf("failed to create Kubernetes Job: %v", err)
    }

    // Wait for the Job to complete
    if err := e.waitForJobCompletion(ctx, jobName, options.Namespace); err != nil {
        // Attempt to fetch logs with retries
        logs, logErr := e.FetchLogsWithRetry(ctx, jobName, options.Namespace, 3, 5*time.Second)
        if logErr != nil {
            return fmt.Errorf("deployment failed: %v; additionally, failed to fetch logs: %v", err, logErr)
        }
        fmt.Printf("Deployment failed. Logs:\n%s\n", logs)
        return fmt.Errorf("deployment failed: %v", err)
    }

    // Retrieve Terraform outputs
    outputs, err := e.retrieveTerraformOutputs(ctx, options)
    if err != nil {
        return fmt.Errorf("failed to retrieve Terraform outputs: %v", err)
    }
    fmt.Printf("Terraform Outputs:\n%v\n", outputs)

    // Clean up Kubernetes resources
    if err := e.cleanupJobAndConfigMap(ctx, jobName, configMapName, secretName, options.Namespace); err != nil {
        return fmt.Errorf("cleanup failed: %v", err)
    }

    return nil
}

// Define PV and PVC as previously described
func definePersistentVolume() *corev1.PersistentVolume {
    pv := &corev1.PersistentVolume{
        ObjectMeta: metav1.ObjectMeta{
            Name: "terraform-state-pv",
        },
        Spec: corev1.PersistentVolumeSpec{
            Capacity: corev1.ResourceList{
                corev1.ResourceStorage: resource.MustParse("1Gi"),
            },
            AccessModes: []corev1.PersistentVolumeAccessMode{
                corev1.ReadWriteOnce,
            },
            PersistentVolumeReclaimPolicy: corev1.PersistentVolumeReclaimRetain,
            StorageClassName:             "manual",
            HostPath: &corev1.HostPathVolumeSource{
                Path: "/mnt/data", // Ensure this path exists on all nodes
            },
        },
    }
    return pv
}

func definePersistentVolumeClaim(namespace string) *corev1.PersistentVolumeClaim {
    pvc := &corev1.PersistentVolumeClaim{
        ObjectMeta: metav1.ObjectMeta{
            Name:      "terraform-state-pvc",
            Namespace: namespace,
        },
        Spec: corev1.PersistentVolumeClaimSpec{
            AccessModes: []corev1.PersistentVolumeAccessMode{
                corev1.ReadWriteOnce,
            },
            StorageClassName: func() *string {
                sc := "manual"
                return &sc
            }(),
            Resources: corev1.ResourceRequirements{
                Requests: corev1.ResourceList{
                    corev1.ResourceStorage: resource.MustParse("1Gi"),
                },
            },
        },
    }
    return pvc
}

func createPersistentVolume(ctx context.Context, clientset *kubernetes.Clientset, pv *corev1.PersistentVolume) error {
    existingPV, err := clientset.CoreV1().PersistentVolumes().Get(ctx, pv.Name, metav1.GetOptions{})
    if err != nil {
        if errors.IsNotFound(err) {
            // Create PV
            _, err = clientset.CoreV1().PersistentVolumes().Create(ctx, pv, metav1.CreateOptions{})
            if err != nil {
                return fmt.Errorf("failed to create PersistentVolume: %v", err)
            }
            fmt.Printf("PersistentVolume %s created.\n", pv.Name)
            return nil
        }
        return fmt.Errorf("failed to get PersistentVolume: %v", err)
    }

    // PV already exists
    fmt.Printf("PersistentVolume %s already exists.\n", existingPV.Name)
    return nil
}

func createPersistentVolumeClaim(ctx context.Context, clientset *kubernetes.Clientset, pvc *corev1.PersistentVolumeClaim) error {
    existingPVC, err := clientset.CoreV1().PersistentVolumeClaims(pvc.Namespace).Get(ctx, pvc.Name, metav1.GetOptions{})
    if err != nil {
        if errors.IsNotFound(err) {
            // Create PVC
            _, err = clientset.CoreV1().PersistentVolumeClaims(pvc.Namespace).Create(ctx, pvc, metav1.CreateOptions{})
            if err != nil {
                return fmt.Errorf("failed to create PersistentVolumeClaim: %v", err)
            }
            fmt.Printf("PersistentVolumeClaim %s created in namespace %s.\n", pvc.Name, pvc.Namespace)
            return nil
        }
        return fmt.Errorf("failed to get PersistentVolumeClaim: %v", err)
    }

    // PVC already exists
    fmt.Printf("PersistentVolumeClaim %s already exists in namespace %s.\n", existingPVC.Name, existingPVC.Namespace)
    return nil
}
