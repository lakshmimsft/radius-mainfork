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

// Deploy triggers a Kubernetes Job to run the Terraform commands for deploying the infrastructure.
func (e *executor) Deploy(ctx context.Context, options Options) error {
    // Validate options.
    if options.WorkingDir == "" || options.RecipeName == "" || options.Namespace == "" {
        return fmt.Errorf("missing required options")
    }

    // Check if the working directory exists.
    if _, err := os.Stat(options.WorkingDir); os.IsNotExist(err) {
        return fmt.Errorf("working directory does not exist")
    }

    // Create a ConfigMap with the Terraform configuration.
    configMapName := fmt.Sprintf("terraform-config-%s", options.RecipeName)
    if err := e.createConfigMap(ctx, configMapName, options); err != nil {
        return fmt.Errorf("failed to create ConfigMap: %v", err)
    }

    // Create a Kubernetes Job to run the Terraform apply command.
    jobName := fmt.Sprintf("terraform-apply-job-%s", options.RecipeName)
    if err := e.createTerraformJob(ctx, jobName, configMapName, "apply", options); err != nil {
        return fmt.Errorf("failed to create Kubernetes Job: %v", err)
    }

    // Wait for the Job to complete.
    if err := e.waitForJobCompletion(ctx, jobName, options.Namespace); err != nil {
        return fmt.Errorf("deployment failed: %v", err)
    }

    // Clean up resources.
    //if err := e.cleanupJobAndConfigMap(ctx, jobName, configMapName, options.Namespace); err != nil {
    //    return fmt.Errorf("cleanup failed: %v", err)
    }

    return nil
}

// Delete triggers a Kubernetes Job to run the Terraform commands for destroying the infrastructure.
func (e *executor) Delete(ctx context.Context, options Options) error {
    // Validate options.
    if options.WorkingDir == "" || options.RecipeName == "" || options.Namespace == "" {
        return fmt.Errorf("missing required options")
    }

    // Check if the working directory exists.
    if _, err := os.Stat(options.WorkingDir); os.IsNotExist(err) {
        return fmt.Errorf("working directory does not exist")
    }

    // Create a ConfigMap with the Terraform configuration.
    configMapName := fmt.Sprintf("terraform-config-%s", options.RecipeName)
    if err := e.createConfigMap(ctx, configMapName, options); err != nil {
        return fmt.Errorf("failed to create ConfigMap: %v", err)
    }

    // Create a Kubernetes Job to run the Terraform destroy command.
    jobName := fmt.Sprintf("terraform-destroy-job-%s", options.RecipeName)
    if err := e.createTerraformJob(ctx, jobName, configMapName, "destroy", options); err != nil {
        return fmt.Errorf("failed to create Kubernetes Job: %v", err)
    }

    // Wait for the Job to complete.
    if err := e.waitForJobCompletion(ctx, jobName, options.Namespace); err != nil {
        return fmt.Errorf("deletion failed: %v", err)
    }

    // Clean up resources.
    if err := e.cleanupJobAndConfigMap(ctx, jobName, configMapName, options.Namespace); err != nil {
        return fmt.Errorf("cleanup failed: %v", err)
    }

    return nil
}

// getKubernetesClient initializes and returns a Kubernetes clientset.
//func getKubernetesClient() (*kubernetes.Clientset, error) {
//    // Try in-cluster config.
//    config, err := rest.InClusterConfig()
//    if err != nil {
//        // If not running in a cluster, use the default config (for local testing).
//        kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
//        config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
//        if err != nil {
//            return nil, fmt.Errorf("failed to create Kubernetes client configuration: %v", err)
//        }
//    }
//
//    // Create the clientset.
//    clientset, err := kubernetes.NewForConfig(config)
//    if err != nil {
//        return nil, fmt.Errorf("failed to create Kubernetes clientset: %v", err)
//    }
//
//    return clientset, nil
//}

// createConfigMap creates a ConfigMap containing the Terraform configuration files.
func (e *executor) createConfigMap(ctx context.Context, configMapName string, options Options) error {
    // Read all files in the working directory.
    files, err := os.ReadDir(options.WorkingDir)
    if err != nil {
        return fmt.Errorf("failed to read working directory: %v", err)
    }

    data := make(map[string]string)
    for _, file := range files {
        if file.IsDir() {
            continue
        }
        filePath := filepath.Join(options.WorkingDir, file.Name())
        content, err := os.ReadFile(filePath)
        if err != nil {
            return fmt.Errorf("failed to read file %s: %v", filePath, err)
        }
        data[file.Name()] = string(content)
    }

    configMap := &corev1.ConfigMap{
        ObjectMeta: metav1.ObjectMeta{
            Name:      configMapName,
            Namespace: options.Namespace,
        },
        Data: data,
    }

    // Create the ConfigMap in Kubernetes.
    _, err = e.clientset.CoreV1().ConfigMaps(options.Namespace).Create(ctx, configMap, metav1.CreateOptions{})
    if err != nil {
        return fmt.Errorf("failed to create ConfigMap: %v", err)
    }

    return nil
} 

// createTerraformJob creates a Kubernetes Job that runs Terraform commands.
func (e *executor) createTerraformJob(ctx context.Context, jobName, configMapName, action string, options Options) error {
    // action can be "apply" or "destroy"
    command := fmt.Sprintf("terraform init && terraform %s -auto-approve", action)

    job := &batchv1.Job{
        ObjectMeta: metav1.ObjectMeta{
            Name:      jobName,
            Namespace: options.Namespace,
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
                                command,
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

    // Create the Job in Kubernetes.
    _, err := e.clientset.BatchV1().Jobs(options.Namespace).Create(ctx, job, metav1.CreateOptions{})
    if err != nil {
        return fmt.Errorf("failed to create Kubernetes Job: %v", err)
    }

    return nil
}

// WaitForJobCompletion waits for the Kubernetes Job to complete. - Add Logger
func (e *executor) waitForJobCompletion(ctx context.Context, jobName, namespace string) error {
    watchInterface, err := e.clientset.BatchV1().Jobs(namespace).Watch(ctx, metav1.ListOptions{
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
    

func (e *executor) cleanupJobAndConfigMap(ctx context.Context, jobName, configMapName, namespace string) error {
    // Delete the Job
    err := e.clientset.BatchV1().Jobs(namespace).Delete(ctx, jobName, metav1.DeleteOptions{})
    if err != nil {
        return fmt.Errorf("failed to delete Job: %v", err)
    }

    // Delete the ConfigMap
    err = e.clientset.CoreV1().ConfigMaps(namespace).Delete(ctx, configMapName, metav1.DeleteOptions{})
    if err != nil {
        return fmt.Errorf("failed to delete ConfigMap: %v", err)
    }

    return nil
}
    
/*
func NewExecutor(clientset *kubernetes.Clientset) *executor {
    return &executor{
        clientset: clientset,
        // Initialize other fields if necessary.
    }
}

// driver
type TerraformDriver struct {
    clientset *kubernetes.Clientset
}

func NewTerraformDriver(clientset *kubernetes.Clientset) *TerraformDriver {
    return &TerraformDriver{
        clientset: clientset,
    }
}

func (d *TerraformDriver) Execute(ctx context.Context, cfg recipes.RecipeConfig) error {
    exec := terraform.NewExecutor(d.clientset)

    options := terraform.Options{
        WorkingDir: cfg.WorkingDir,
        RecipeName: cfg.RecipeName,
        Namespace:  cfg.Namespace, // Assuming Namespace is part of RecipeConfig
        // Set other options as necessary.
    }

    return exec.Deploy(ctx, options)
}

func (d *TerraformDriver) Destroy(ctx context.Context, cfg recipes.RecipeConfig) error {
    exec := terraform.NewExecutor(d.clientset)

    options := terraform.Options{
        WorkingDir: cfg.WorkingDir,
        RecipeName: cfg.RecipeName,
        Namespace:  cfg.Namespace,
        // Set other options as necessary.
    }

    return exec.Delete(ctx, options)
}

// engine
func RunRecipe(ctx context.Context, cfg recipes.RecipeConfig) error {
    // Initialize the Kubernetes client.
    kubeClient, err := kubernetes.NewClient()
    if err != nil {
        return err
    }

    // Initialize the TerraformDriver.
    terraformDriver := driver.NewTerraformDriver(kubeClient.Clientset)

    // Execute the recipe.
    return terraformDriver.Execute(ctx, cfg)
}

*/