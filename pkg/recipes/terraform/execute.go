/*
Copyright 2023 The Radius Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package terraform
package terraform

import (
    "context"
    "errors"
    "fmt"
    "os"
    "path/filepath"

    batchv1 "k8s.io/api/batch/v1"
    corev1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/rest"
    "k8s.io/client-go/tools/clientcmd"
)

func Deploy(ctx context.Context, workingDir, recipeName string) error {
    // Check if the working directory exists.
    if _, err := os.Stat(workingDir); os.IsNotExist(err) {
        return fmt.Errorf("working directory does not exist")
    }

    // Initialize the Kubernetes client.
    clientset, err := getKubernetesClient()
    if err != nil {
        return fmt.Errorf("failed to initialize Kubernetes client: %v", err)
    }

    // Create a ConfigMap with the Terraform configuration.
    configMapName := fmt.Sprintf("terraform-config-%s", recipeName)
    if err := createConfigMap(ctx, clientset, configMapName, workingDir); err != nil {
        return fmt.Errorf("failed to create ConfigMap: %v", err)
    }

    // Create a Kubernetes Job to run the Terraform apply command.
    jobName := fmt.Sprintf("terraform-apply-job-%s", recipeName)
    if err := createTerraformApplyJob(ctx, clientset, jobName, configMapName); err != nil {
        return fmt.Errorf("failed to create Kubernetes Job: %v", err)
    }

    // Wait for the Job to complete.
    err = WaitForJobCompletion(ctx, clientset, jobName, "default")
    if err != nil {
        return fmt.Errorf("deployment failed: %v", err)
    }

    // clean up resources - keep for debugging
    // err = cleanupJobAndConfigMap(ctx, clientset, jobName, configMapName, "default")
    // if err != nil {
    //     return fmt.Errorf("cleanup failed: %v", err)
    // }

    return nil
}

// Delete triggers a Kubernetes Job to run the Terraform commands for destroying the infrastructure.
func Delete(ctx context.Context, workingDir, recipeName string) error {
    // Check if the working directory exists.
    if _, err := os.Stat(workingDir); os.IsNotExist(err) {
        return fmt.Errorf("working directory does not exist")
    }

    // Initialize the Kubernetes client.
    clientset, err := getKubernetesClient()
    if err != nil {
        return fmt.Errorf("failed to initialize Kubernetes client: %v", err)
    }

    // Create a ConfigMap with the Terraform configuration.
    configMapName := fmt.Sprintf("terraform-config-%s", recipeName)
    if err := createConfigMap(ctx, clientset, configMapName, workingDir); err != nil {
        return fmt.Errorf("failed to create ConfigMap: %v", err)
    }

    // Create a Kubernetes Job to run the Terraform destroy command.
    jobName := fmt.Sprintf("terraform-destroy-job-%s", recipeName)
    if err := createTerraformDestroyJob(ctx, clientset, jobName, configMapName); err != nil {
        return fmt.Errorf("failed to create Kubernetes Job: %v", err)
    }

    // Wait for the Job to complete.
    err = WaitForJobCompletion(ctx, clientset, jobName, "default")
    if err != nil {
        return fmt.Errorf("deletion failed: %v", err)
    }

    // Clean up resources.
    err = cleanupJobAndConfigMap(ctx, clientset, jobName, configMapName, "default")
    if err != nil {
        return fmt.Errorf("cleanup failed: %v", err)
    }

    return nil
}

// createTerraformDestroyJob creates a Kubernetes Job that runs 'terraform destroy'.
func createTerraformDestroyJob(ctx context.Context, clientset *kubernetes.Clientset, jobName, configMapName string) error {
    job := &batchv1.Job{
        ObjectMeta: metav1.ObjectMeta{
            Name: jobName,
        },
        Spec: batchv1.JobSpec{
            Template: corev1.PodTemplateSpec{
                ObjectMeta: metav1.ObjectMeta{
                    Labels: map[string]string{
                        "job-name": jobName,
                    },
                },
                Spec: corev1.PodSpec{
                    RestartPolicy: corev1.RestartPolicyNever,
                    Containers: []corev1.Container{
                        {
                            Name:  "terraform",
                            Image: "hashicorp/terraform:light",
                            Command: []string{
                                "sh",
                                "-c",
                                "terraform init && terraform destroy -auto-approve",
                            },
                            WorkingDir: "/app/terraform",
                            VolumeMounts: []corev1.VolumeMount{
                                {
                                    Name:      "terraform-config",
                                    MountPath: "/app/terraform",
                                },
                            },
                        },
                    },
                    Volumes: []corev1.Volume{
                        {
                            Name: "terraform-config",
                            VolumeSource: corev1.VolumeSource{
                                ConfigMap: &corev1.ConfigMapVolumeSource{
                                    LocalObjectReference: corev1.LocalObjectReference{
                                        Name: configMapName,
                                    },
                                },
                            },
                        },
                    },
                },
            },
        },
    }

    // Create the Job in the 'default' namespace - temporary.
    _, err := clientset.BatchV1().Jobs("default").Create(ctx, job, metav1.CreateOptions{})
    if err != nil {
        return fmt.Errorf("failed to create Kubernetes Job: %v", err)
    }

    return nil
}

// createTerraformApplyJob creates a Kubernetes Job that runs 'terraform apply'.
func createTerraformApplyJob(ctx context.Context, clientset *kubernetes.Clientset, jobName, configMapName string) error {
    job := &batchv1.Job{
        ObjectMeta: metav1.ObjectMeta{
            Name: jobName,
        },
        Spec: batchv1.JobSpec{
            Template: corev1.PodTemplateSpec{
                ObjectMeta: metav1.ObjectMeta{
                    Labels: map[string]string{
                        "job-name": jobName,
                    },
                },
                Spec: corev1.PodSpec{
                    RestartPolicy: corev1.RestartPolicyNever,
                    Containers: []corev1.Container{
                        {
                            Name:  "terraform",
                            Image: "hashicorp/terraform:light",
                            Command: []string{
                                "sh",
                                "-c",
                                "terraform init && terraform apply -auto-approve",
                            },
                            WorkingDir: "/app/terraform",
                            VolumeMounts: []corev1.VolumeMount{
                                {
                                    Name:      "terraform-config",
                                    MountPath: "/app/terraform",
                                },
                            },
                        },
                    },
                    Volumes: []corev1.Volume{
                        {
                            Name: "terraform-config",
                            VolumeSource: corev1.VolumeSource{
                                ConfigMap: &corev1.ConfigMapVolumeSource{
                                    LocalObjectReference: corev1.LocalObjectReference{
                                        Name: configMapName,
                                    },
                                },
                            },
                        },
                    },
                },
            },
        },
    }

    // Create the Job in the 'default' namespace - temporary.
    _, err := clientset.BatchV1().Jobs("default").Create(ctx, job, metav1.CreateOptions{})
    if err != nil {
        return fmt.Errorf("failed to create Kubernetes Job: %v", err)
    }

    return nil
}

// getKubernetesClient initializes and returns a Kubernetes clientset.
func getKubernetesClient() (*kubernetes.Clientset, error) {
    // Try in-cluster config.
    config, err := rest.InClusterConfig()
    if err != nil {
        // If not running in a cluster, use the default config (for local testing).
        kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
        config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
        if err != nil {
            return nil, fmt.Errorf("failed to create Kubernetes client configuration: %v", err)
        }
    }

    // Create the clientset.
    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        return nil, fmt.Errorf("failed to create Kubernetes clientset: %v", err)
    }

    return clientset, nil
}

// createConfigMap creates a ConfigMap containing the Terraform configuration files.
func createConfigMap(ctx context.Context, clientset *kubernetes.Clientset, configMapName, workingDir string) error {
    // Read all files in the working directory.
    files, err := os.ReadDir(workingDir)
    if err != nil {
        return fmt.Errorf("failed to read working directory: %v", err)
    }

    data := make(map[string]string)
    for _, file := range files {
        if file.IsDir() {
            continue
        }
        filePath := filepath.Join(workingDir, file.Name())
        content, err := os.ReadFile(filePath)
        if err != nil {
            return fmt.Errorf("failed to read file %s: %v", filePath, err)
        }
        data[file.Name()] = string(content)
    }

    configMap := &corev1.ConfigMap{
        ObjectMeta: metav1.ObjectMeta{
            Name: configMapName,
        },
        Data: data,
    }

    // Create the ConfigMap in the 'default' namespace.
    _, err = clientset.CoreV1().ConfigMaps("default").Create(ctx, configMap, metav1.CreateOptions{})
    if err != nil {
        return fmt.Errorf("failed to create ConfigMap: %v", err)
    }

    return nil
}

// createTerraformJob creates a Kubernetes Job that runs Terraform commands.
func createTerraformJob(ctx context.Context, clientset *kubernetes.Clientset, jobName, configMapName string) error {
    job := &batchv1.Job{
        ObjectMeta: metav1.ObjectMeta{
            Name: jobName,
        },
        Spec: batchv1.JobSpec{
            Template: corev1.PodTemplateSpec{
                ObjectMeta: metav1.ObjectMeta{
                    Labels: map[string]string{
                        "job-name": jobName,
                    },
                },
                Spec: corev1.PodSpec{
                    RestartPolicy: corev1.RestartPolicyNever,
                    Containers: []corev1.Container{
                        {
                            Name:  "terraform",
                            Image: "hashicorp/terraform:light",
                            Command: []string{
                                "sh",
                                "-c",
                                "terraform init && terraform apply -auto-approve",
                            },
                            WorkingDir: "/app/terraform",
                            VolumeMounts: []corev1.VolumeMount{
                                {
                                    Name:      "terraform-config",
                                    MountPath: "/app/terraform",
                                },
                            },
                        },
                    },
                    Volumes: []corev1.Volume{
                        {
                            Name: "terraform-config",
                            VolumeSource: corev1.VolumeSource{
                                ConfigMap: &corev1.ConfigMapVolumeSource{
                                    LocalObjectReference: corev1.LocalObjectReference{
                                        Name: configMapName,
                                    },
                                },
                            },
                        },
                    },
                },
            },
        },
    }

    // Create the Job in the 'default' namespace.
    _, err := clientset.BatchV1().Jobs("default").Create(ctx, job, metav1.CreateOptions{})
    if err != nil {
        return fmt.Errorf("failed to create Kubernetes Job: %v", err)
    }

    // Wait for the Job to complete
    err = WaitForJobCompletion(ctx, clientset, jobName, "default")
    if err != nil {
        return fmt.Errorf("deployment failed: %v", err)
    }

    // err = cleanupJobAndConfigMap(ctx, clientset, jobName, configMapName, "default")
    // if err != nil {
    //     return fmt.Errorf("cleanup failed: %v", err)
    // }

    return nil
}

// WaitForJobCompletion waits for the Kubernetes Job to complete. - Add Logger
func WaitForJobCompletion(ctx context.Context, clientset *kubernetes.Clientset, jobName string, namespace string) error {
    watchInterface, err := clientset.BatchV1().Jobs(namespace).Watch(ctx, metav1.ListOptions{
        FieldSelector: fmt.Sprintf("metadata.name=%s", jobName),
    })
    if err != nil {
        return fmt.Errorf("failed to start watch on Job: %v", err)
    }
    defer watchInterface.Stop()

    for {
        select {
        case <-ctx.Done():
            return fmt.Errorf("context canceled or timed out")
        case event, ok := <-watchInterface.ResultChan():
            if !ok {
                return fmt.Errorf("watch channel closed")
            }
            job, ok := event.Object.(*batchv1.Job)
            if !ok {
                return fmt.Errorf("unexpected type received from watch")
            }

            if job.Status.Succeeded > 0 {
                fmt.Println("Job completed successfully")
                return nil
            }

            if job.Status.Failed > 0 {
                return fmt.Errorf("Job failed")
            }
        }
    }
}

func cleanupJobAndConfigMap(ctx context.Context, clientset *kubernetes.Clientset, jobName, configMapName, namespace string) error {
    // Delete the Job
    err := clientset.BatchV1().Jobs(namespace).Delete(ctx, jobName, metav1.DeleteOptions{})
    if err != nil {
        return fmt.Errorf("failed to delete Job: %v", err)
    }

    // Delete the ConfigMap
    err = clientset.CoreV1().ConfigMaps(namespace).Delete(ctx, configMapName, metav1.DeleteOptions{})
    if err != nil {
        return fmt.Errorf("failed to delete ConfigMap: %v", err)
    }

    return nil
}
